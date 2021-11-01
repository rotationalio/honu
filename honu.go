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
func Open(uri string, conf config.ReplicaConfig) (db *DB, err error) {
	var dsn *DSN
	if dsn, err = ParseDSN(uri); err != nil {
		return nil, err
	}

	db = &DB{}
	if db.vm, err = NewVersionManager(conf); err != nil {
		return nil, err
	}

	switch dsn.Scheme {
	case "leveldb":
		// TODO: allow leveldb options to be passed to OpenFile
		// TODO: multiple leveldb databases for different namespaces
		if db.engine, err = leveldb.Open(dsn.Path, conf); err != nil {
			return nil, err
		}
	case "badger", "badgerdb":
		if db.engine, err = badger.Open(conf); err != nil {
			return nil, err
		}
	case "pebble", "pebbledb":
		if db.engine, err = pebble.Open(conf); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unhandled database scheme %q", dsn.Scheme)
	}

	return db, nil
}

// Close the database, allowing no further interactions.
func (d *DB) Close() error {
	return d.engine.Close()
}

// Get the latest version of the object stored by the key.
// TODO: provide read options to the underlying database.
func (d *DB) Get(key []byte, options string) (value []byte, err error) {
	// TODO: refactor this into an options slice for faster checking
	store, ok := d.engine.(engine.Store)
	if !ok {
		return nil, errors.New("underlying engine doesn't support Get accesses")
	}

	// Fetch the value from the database
	if value, err = store.Get(key, options); err != nil {
		// TODO: wrap the engine error in standard honu errors?
		return nil, err
	}

	// Parse the object record to extract the value data
	obj := new(pb.Object)
	if err = proto.Unmarshal(value, obj); err != nil {
		// TODO: better error message here
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

// Put a new value to the specified key and update the version.
// TODO: provide write options to the underlying database.
func (d *DB) Put(key, value []byte, options string) (err error) {
	// TODO: refactor this into an options slice for faster checking
	store, ok := d.engine.(engine.Store)
	if !ok {
		return errors.New("underlying engine doesn't support Put accesses")
	}

	// Get or Create the previous version
	var data []byte
	var obj *pb.Object
	if data, err = store.Get(key, options); err != nil {
		if errors.Is(err, engine.ErrNotFound) {
			obj = &pb.Object{
				Key:       key,
				Namespace: "default",
			}
		} else {
			return err
		}
	} else {
		obj = new(pb.Object)
		if err = proto.Unmarshal(data, obj); err != nil {
			return err
		}
	}

	// Update the version with the new data
	obj.Data = value
	if err = d.vm.Update(obj); err != nil {
		return err
	}

	// Put the version back onto disk
	if data, err = proto.Marshal(obj); err != nil {
		return err
	}
	if err = store.Put(key, data, options); err != nil {
		return err
	}

	return nil
}

// Delete the object represented by the key, creating a tombstone object.
// TODO: provide write options to the underlying database.
func (d *DB) Delete(key []byte, options string) (err error) {
	// TODO: refactor this into an options slice for faster checking
	store, ok := d.engine.(engine.Store)
	if !ok {
		return errors.New("underlying engine doesn't support Delete accesses")
	}

	var data []byte
	if data, err = store.Get(key, options); err != nil {
		if errors.Is(err, engine.ErrNotFound) {
			return nil
		}
		return err
	}

	// Unmarshal the version information
	obj := new(pb.Object)
	if err = proto.Unmarshal(data, obj); err != nil {
		return err
	}

	// Don't save the data back to disk
	obj.Data = nil

	// Create a tombstone for the data
	if err = d.vm.Delete(obj); err != nil {
		return err
	}

	// Put the version back onto disk
	if data, err = proto.Marshal(obj); err != nil {
		return err
	}
	// TODO allow writeoptions to be passed to this put
	if err = store.Put(key, data, ""); err != nil {
		return err
	}
	return nil
}

// Iter over a subset of keys specified by the prefix.
// TODO: provide better mechanisms for iteration.
func (d *DB) Iter(prefix []byte) (i iterator.Iterator, err error) {
	// TODO: refactor this into an options slice for faster checking
	iter, ok := d.engine.(engine.Iterator)
	if !ok {
		return nil, errors.New("underlying engine doesn't support Iter accesses")
	}
	return iter.Iter(prefix)
}

// Object returns metadata associated with the latest object stored by the key.
func (d *DB) Object(key []byte, options string) (_ *pb.Object, err error) {
	// TODO: refactor this into an options slice for faster checking
	store, ok := d.engine.(engine.Store)
	if !ok {
		return nil, errors.New("underlying engine doesn't support Object accesses")
	}

	// Fetch the value from the database
	var value []byte
	if value, err = store.Get(key, options); err != nil {
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
