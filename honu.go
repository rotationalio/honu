/*
Package honu provides a thin wrapper over an embedded database (leveldb, sqlite) that
provides version history to object changes and anti-entropy replication.
*/
package honu

import (
	"errors"
	"fmt"

	"github.com/rotationalio/honu/config"
	engine "github.com/rotationalio/honu/engines"
	"github.com/rotationalio/honu/engines/badger"
	"github.com/rotationalio/honu/engines/leveldb"
	"github.com/rotationalio/honu/engines/pebble"
	"github.com/rotationalio/honu/iterator"
	pb "github.com/rotationalio/honu/object"
	opts "github.com/rotationalio/honu/options"
	"google.golang.org/protobuf/proto"
)

// DB is a Honu embedded database.
// Currently DB simply wraps a leveldb database
type DB struct {
	engine engine.Engine
	vm     *VersionManager
}

// Open a replicated embedded database with the specified URI. Database URIs should
// specify protocol:///relative/path/to/db for embedded databases. For absolute paths,
// specify protocol:////absolute/path/to/db.
func Open(uri string, options ...config.Option) (db *DB, err error) {
	// Create a configuration from the options passed in.
	var conf config.Config
	if conf, err = config.New(options...); err != nil {
		return nil, err
	}

	var dsn *DSN
	if dsn, err = ParseDSN(uri); err != nil {
		return nil, err
	}

	db = &DB{}
	if db.vm, err = NewVersionManager(conf.Versions); err != nil {
		return nil, err
	}

	switch dsn.Scheme {
	case "leveldb":
		// TODO: multiple leveldb databases for different namespaces
		if db.engine, err = leveldb.Open(dsn.Path, conf); err != nil {
			return nil, err
		}
	case "badger", "badgerdb":
		if db.engine, err = badger.Open(conf); err != nil {
			return nil, err
		}
	case "pebble", "pebbledb":
		if db.engine, err = pebble.Open(dsn.Path, conf); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unhandled database scheme %q", dsn.Scheme)
	}

	return db, nil
}

// Close the database, allowing no further interactions.
func (db *DB) Close() error {
	return db.engine.Close()
}

// Object returns metadata associated with the latest object stored by the key.
// Object is the Get function to use if you want to fetch tombstones, otherwise use Get
// which will return a not found error.
func (db *DB) Object(key []byte, options ...opts.Option) (_ *pb.Object, err error) {
	var tx engine.Transaction
	if tx, err = db.engine.Begin(true); err != nil {
		return nil, err
	}
	defer tx.Finish()

	// Collect the options
	var cfg *opts.Options
	if cfg, err = opts.New(options...); err != nil {
		return nil, err
	}

	// Fetch the value from the database
	var value []byte
	if value, err = tx.Get(key, cfg); err != nil {
		// TODO: should we wrap the leveldb error?
		return nil, err
	}

	// Parse the object record to extract the value data
	obj := new(pb.Object)
	if err = proto.Unmarshal(value, obj); err != nil {
		// TODO: better error message here
		return nil, err
	}
	return obj, nil
}

// Get the latest version of the object stored by the key.
func (db *DB) Get(key []byte, options ...opts.Option) (value []byte, err error) {
	var obj *pb.Object
	if obj, err = db.Object(key, options...); err != nil {
		return nil, err
	}

	if obj.Tombstone() {
		// The object is deleted, so return not found
		// TODO: standardize error messages
		return nil, engine.ErrNotFound
	}

	// Return the wrapped data
	return obj.Data, nil
}

// UpdateType is an intermediate solution to tracking what's happening to the version
// history when direct modifications are applied to the database.
// NOTE: this data type is subject to change in later versions and should be treated
// as a prototype only in production code.
type UpdateType uint8

const (
	UpdateNoChange UpdateType = iota // No change occurred (nothing was written to disk)
	UpdateForced                     // The update was forced, so the previous version was not checked
	UpdateLinear                     // The previous version is the parent of the updating version
	UpdateStomp                      // The previous version is concurrent but has a lower precedence than the updating version
	UpdateSkip                       // The previous version is later but is not concurrent nor linear from the updating version

)

func (u UpdateType) String() string {
	switch u {
	case UpdateNoChange:
		return "unchanged"
	case UpdateForced:
		return "forced"
	case UpdateLinear:
		return "linear"
	case UpdateStomp:
		return "stomped"
	case UpdateSkip:
		return "skipped"
	default:
		return "unknown"
	}
}

// Update an object directly in the database without modifying its version information.
// Update is to Put as Object is to Get - use Update when manually modifying the data
// store, for example during replication, but not for normal DB operations. Update also
// returns the type of update that ocurred, relative to the previous version.
func (db *DB) Update(obj *pb.Object, options ...opts.Option) (update UpdateType, err error) {
	var tx engine.Transaction
	if tx, err = db.engine.Begin(false); err != nil {
		return UpdateNoChange, err
	}
	defer tx.Finish()

	// Collect the options
	var cfg *opts.Options
	if cfg, err = opts.New(options...); err != nil {
		return UpdateNoChange, err
	}

	// If the default namespace is specified use the object's namespace to ensure that
	// if the user did not supply Namespace, Update still works. If the object was
	// already in the default namespace, this should not cause a change to happen. There
	// is an edge case where the user supplies options.WithNamespace("default") and an
	// object that is not in the default namespace and the user option will be ignored
	// in favor of the object's original namespace; but this is an unlikely case.
	if cfg.Namespace == opts.NamespaceDefault {
		cfg.Namespace = obj.Namespace
	}

	if !cfg.Force {
		// Check the namespace and that it matches the object
		if cfg.Namespace != obj.Namespace {
			return UpdateNoChange, errors.New("options namespace does not match object namespace")
		}

		// Check that the version is later than the version being written to disk
		var (
			prevData []byte
			prev     = new(pb.Object)
		)
		if prevData, err = tx.Get(obj.Key, cfg); err != nil {
			if !errors.Is(err, engine.ErrNotFound) {
				return UpdateNoChange, fmt.Errorf("could not check previous version: %v", err)
			}
		} else {
			if err = proto.Unmarshal(prevData, prev); err != nil {
				return UpdateNoChange, fmt.Errorf("could not unmarshal previous version: %v", err)
			}
		}

		if !obj.Version.IsLater(prev.Version) {
			return UpdateNoChange, fmt.Errorf("cannot update object, it is not a later version then the current object")
		}

		// Determine the update type based on the previous version
		// NOTE: all update conditions imply the current version is later than previous
		switch {
		case obj.Version.Stomps(prev.Version):
			update = UpdateStomp
		case obj.Version.Skips(prev.Version):
			update = UpdateSkip
		case obj.Version.LinearFrom(prev.Version):
			update = UpdateLinear
		default:
			return UpdateNoChange, fmt.Errorf("cannot determine update relationship in the version history")
		}

	} else {
		// Report that the update was forced
		update = UpdateForced
	}

	// Put the version directly to disk
	var data []byte
	if data, err = proto.Marshal(obj); err != nil {
		return UpdateNoChange, err
	}

	if err = tx.Put(obj.Key, data, cfg); err != nil {
		return UpdateNoChange, err
	}
	return update, nil
}

// Put a new value to the specified key and update the version.
func (db *DB) Put(key, value []byte, options ...opts.Option) (_ *pb.Object, err error) {
	var tx engine.Transaction
	if tx, err = db.engine.Begin(false); err != nil {
		return nil, err
	}
	defer tx.Finish()

	// Collect the options
	var cfg *opts.Options
	if cfg, err = opts.New(options...); err != nil {
		return nil, err
	}

	// Get or Create the previous version
	var data []byte
	var obj *pb.Object
	if data, err = tx.Get(key, cfg); err != nil {
		if errors.Is(err, engine.ErrNotFound) {
			obj = &pb.Object{
				Key:       key,
				Namespace: cfg.Namespace,
			}
		} else {
			return nil, err
		}
	} else {
		obj = new(pb.Object)
		if err = proto.Unmarshal(data, obj); err != nil {
			return nil, err
		}
	}

	// Update the version with the new data
	obj.Data = value
	if err = db.vm.Update(obj); err != nil {
		return nil, err
	}

	// Put the version back onto disk
	if data, err = proto.Marshal(obj); err != nil {
		return nil, err
	}
	if err = tx.Put(key, data, cfg); err != nil {
		return nil, err
	}

	// Test to make sure obj.Data is not modified if value is modified.
	return obj, nil
}

// Delete the object represented by the key, creating a tombstone object.
func (db *DB) Delete(key []byte, options ...opts.Option) (_ *pb.Object, err error) {
	var tx engine.Transaction
	if tx, err = db.engine.Begin(false); err != nil {
		return nil, err
	}
	defer tx.Finish()

	// Collect the options
	var cfg *opts.Options
	if cfg, err = opts.New(options...); err != nil {
		return nil, err
	}

	var data []byte
	if data, err = tx.Get(key, cfg); err != nil {
		if errors.Is(err, engine.ErrNotFound) {
			return nil, err
		}
		return nil, err
	}

	// Unmarshal the version information
	obj := new(pb.Object)
	if err = proto.Unmarshal(data, obj); err != nil {
		return nil, err
	}

	// Don't save the data back to disk
	obj.Data = nil

	// Create a tombstone for the data
	if err = db.vm.Delete(obj); err != nil {
		return nil, err
	}

	// Put the version back onto disk
	if data, err = proto.Marshal(obj); err != nil {
		return nil, err
	}

	if err = tx.Put(key, data, cfg); err != nil {
		return nil, err
	}
	return obj, nil
}

// Iter over a subset of keys specified by the prefix.
// TODO: provide better mechanisms for iteration.
func (db *DB) Iter(prefix []byte, options ...opts.Option) (i iterator.Iterator, err error) {
	// Collect the options
	var cfg *opts.Options
	if cfg, err = opts.New(options...); err != nil {
		return nil, err
	}

	// TODO: refactor this into an options slice for faster checking
	iter, ok := db.engine.(engine.Iterator)
	if !ok {
		return nil, errors.New("underlying engine doesn't support Iter accesses")
	}
	return iter.Iter(prefix, cfg)
}
