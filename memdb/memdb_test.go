package memdb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func Test_DB(t *testing.T) {
	db := New()

	_ = db.Put([]byte("hello"), []byte("world"))
	value, err := db.Get([]byte("hello"))
	fmt.Println(value, err)

	buf := make([][4]byte, 10000)
	for i := range buf {
		binary.LittleEndian.PutUint32(buf[i][:], uint32(i))
	}
	for i := range buf {
		_ = db.Put(buf[i][:], buf[i][:])
	}
	for i := range buf {
		value, _ := db.Get(buf[i][:])
		if !bytes.Equal(value, buf[i][:]) {
			t.Errorf("db.Get(%v) = %v\n", i, value)
		}
	}
}
