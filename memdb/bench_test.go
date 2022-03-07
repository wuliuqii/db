package memdb

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"testing"
)

func BenchmarkPut(b *testing.B) {
	buf := make([][4]byte, b.N)
	for i := range buf {
		binary.LittleEndian.PutUint32(buf[i][:], uint32(i))
	}

	b.ResetTimer()
	db := New()
	for i := range buf {
		_ = db.Put(buf[i][:], buf[i][:])
	}
}

func BenchmarkPutRandom(b *testing.B) {
	buf := make([][4]byte, b.N)
	for i := range buf {
		binary.LittleEndian.PutUint32(buf[i][:], uint32(rand.Int()))
	}

	b.ResetTimer()
	db := New()
	for i := range buf {
		_ = db.Put(buf[i][:], buf[i][:])
	}
}

func BenchmarkGet(b *testing.B) {
	buf := make([][4]byte, b.N)
	for i := range buf {
		binary.LittleEndian.PutUint32(buf[i][:], uint32(i))
	}

	db := New()
	for i := range buf {
		_ = db.Put(buf[i][:], buf[i][:])
	}

	b.ResetTimer()
	for i := range buf {
		value, _ := db.Get(buf[i][:])
		if !bytes.Equal(value, buf[i][:]) {
			b.Errorf("db.Get(%v) = %v\n", i, value)
		}
	}
}

func BenchmarkGetRandom(b *testing.B) {
	buf := make([][4]byte, b.N)
	for i := range buf {
		binary.LittleEndian.PutUint32(buf[i][:], uint32(rand.Int()))
	}

	db := New()
	for i := range buf {
		_ = db.Put(buf[i][:], buf[i][:])
	}

	b.ResetTimer()
	for i := range buf {
		value, _ := db.Get(buf[i][:])
		if !bytes.Equal(value, buf[i][:]) {
			b.Errorf("db.Get(%v) = %v\n", i, value)
		}
	}
}
