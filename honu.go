/*
Package honu provides a thin wrapper over an embedded database (leveldb, sqlite) that
provides version history to object changes and anti-entropy replication.
*/
package honu

import (
	"errors"
	"fmt"

	"github.com/rotationalio/honu/config"
	"github.com/rotationalio/honu/iterator"
	pb "github.com/rotationalio/honu/proto/v1"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"google.golang.org/protobuf/proto"
)

// DB is a Honu embedded database.
// Currently DB simply wraps a leveldb database
type DB struct {
	ldb *leveldb.DB
	vm  *VersionManager
}

// Open a replicated embedded database with the specified URI. Database URIs should
// specify protocol:///relative/path/to/db for embedded databases. For absolute paths,
// specify protocol:////absolute/path/to/db.
func Open(uri string, conf config.ReplicaConfig) (db *DB, err error) {
	var dsn *DSN
	if dsn, err = ParseDSN(uri); err != nil {
		return nil, err
	}

	switch dsn.Scheme {
	case "leveldb":
		// TODO: allow leveldb options to be passed to OpenFile
		// TODO: multiple leveldb databases for different namespaces
		db = &DB{}
		if db.ldb, err = leveldb.OpenFile(dsn.Path, nil); err != nil {
			return nil, err
		}
		if db.vm, err = NewVersionManager(conf); err != nil {
			return nil, err
		}
		return db, nil
	case "sqlite", "sqlite3":
		return nil, errors.New("sqlite support is currently not implemented")
	default:
		return nil, fmt.Errorf("unhandled database scheme %q", dsn.Scheme)
	}
}

// Close the database, allowing no further interactions.
func (d *DB) Close() error {
	return d.ldb.Close()
}

// Get the latest version of the object stored by the key.
// TODO: provide read options to the underlying database.
func (d *DB) Get(key []byte) (value []byte, err error) {
	// Fetch the value from the database
	if value, err = d.ldb.Get(key, nil); err != nil {
		// TODO: should we wrap the leveldb error?
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
		return nil, leveldb.ErrNotFound
	}

	// Return the wrapped data
	return obj.Data, nil

}

// Put a new value to the specified key and update the version.
// TODO: provide write options to the underlying database.
func (d *DB) Put(key, value []byte) (err error) {
	// Get or Create the previous version
	var data []byte
	var obj *pb.Object
	if data, err = d.ldb.Get(key, nil); err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
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
	if err = d.ldb.Put(key, data, nil); err != nil {
		return err
	}

	return nil
}

// Delete the object represented by the key, creating a tombstone object.
// TODO: provide write options to the underlying database.
func (d *DB) Delete(key []byte) (err error) {
	var data []byte
	if data, err = d.ldb.Get(key, nil); err != nil {
		if errors.Is(err, leveldb.ErrNotFound) {
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
	if err = d.ldb.Put(key, data, nil); err != nil {
		return err
	}
	return nil
}

// Iter over a subset of keys specified by the prefix.
// TODO: provide better mechanisms for iteration.
func (d *DB) Iter(prefix []byte) (i iterator.Iterator, err error) {
	var slice *util.Range
	if len(prefix) > 0 {
		slice = util.BytesPrefix(prefix)
	}
	return iterator.NewLevelDBIterator(d.ldb.NewIterator(slice, nil)), nil
}

// Object returns metadata associated with the latest object stored by the key.
func (d *DB) Object(key []byte) (_ *pb.Object, err error) {
	// Fetch the value from the database
	var value []byte
	if value, err = d.ldb.Get(key, nil); err != nil {
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
