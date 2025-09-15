package iterator

import (
	"go.etcd.io/bbolt"
	"go.rtnl.ai/honu/pkg/store/key"
	"go.rtnl.ai/honu/pkg/store/object"
)

func New(cursor *bbolt.Cursor) Iterator {
	return &Cursor{cursor: cursor}
}

// Cursor is a wrapper around a bbolt cursor that implements the Iterator interface.
type Cursor struct {
	cursor  *bbolt.Cursor
	started bool
	key     []byte
	value   []byte
	err     error
}

func (c *Cursor) Key() key.Key {
	if len(c.key) != 0 {
		k := make(key.Key, len(c.key))
		copy(k, c.key)
		return k
	}
	return nil
}

func (c *Cursor) Object() object.Object {
	if len(c.value) != 0 {
		obj := make(object.Object, len(c.value))
		copy(obj, c.value)
		return obj
	}
	return nil
}

func (c *Cursor) Error() error {
	return c.err
}

// Seek to the specified key using binary search. If an exact match is not found, the
// cursor is positioned at the next key that is greater than the given key. If the key
// is at the end of the collection, then false is returned.
func (c *Cursor) Seek(key []byte) bool {
	if c.released() {
		return false
	}

	c.key, c.value = c.cursor.Seek(key)
	c.started = true
	return c.key != nil
}

// Iterate in forward order. If Next() is the first call, the iterator will seek to the
// first item. If Next() returns false then the iterator has been exhausted.
func (c *Cursor) Next() bool {
	if c.released() {
		return false
	}

	if !c.started {
		c.key, c.value = c.cursor.First()
		c.started = true
		return c.key != nil
	}

	c.key, c.value = c.cursor.Next()
	return c.key != nil
}

// Iterate in reverse order. If Prev() is the first call, the iterator will seek to the
// last item, allowing you to easily perform a reverse iteration.
func (c *Cursor) Prev() bool {
	if c.released() {
		return false
	}

	if !c.started {
		c.key, c.value = c.cursor.Last()
		c.started = true
		return c.key != nil
	}

	c.key, c.value = c.cursor.Prev()
	return c.key != nil
}

// Move the cursor to the first item.
func (c *Cursor) First() bool {
	if c.released() {
		return false
	}

	c.key, c.value = c.cursor.First()
	c.started = true
	return c.key != nil
}

// Move the cursor to the last item.
func (c *Cursor) Last() bool {
	if c.released() {
		return false
	}

	c.key, c.value = c.cursor.Last()
	c.started = true
	return c.key != nil
}

func (c *Cursor) Release() {
	c.cursor = nil
}

func (c *Cursor) released() bool {
	return c.cursor == nil
}
