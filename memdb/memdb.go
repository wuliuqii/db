package memdb

// in-memory key-value table
type table interface {
	// put the value for the given key. Tt overwrites any previous value
	// for that key.
	put(key, value []byte) error
	// get the value for the given key. It returns ErrNotFound if not
	// contain the key
	get(key []byte) ([]byte, error)
}

// DB is an in-memory key/value database
type DB struct {
	table table
}

func (db *DB) Put(key, value []byte) error {
	return db.table.put(key, value)
}

func (db *DB) Get(key []byte) ([]byte, error) {
	return db.table.get(key)
}

func New() *DB {
	return &DB{newSkipList(bytesComparator{})}
}
