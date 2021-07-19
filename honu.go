package honu

import "github.com/rotationalio/honu/iterator"

type DB struct {
}

func Open(dsn string) *DB {
	return nil
}

func Get(key []byte) (value []byte, err error) {
	return nil, nil
}

func Put(key, value []byte) (err error) {
	return nil
}

func Delete(key []byte) (err error) {
	return nil
}

func Iter(prefix []byte) (i *iterator.Iterator, err error) {
	return nil, nil
}

func (db *DB) Version(key []byte) (o *v1.Object, err error) {
	return nil, nil
}
