package kiwi

import (
	"encoding"
	"io"
)

// Client is the full key-value interface to a data store
type Client interface {
	io.Closer

	Crud
	Iterator
}

// Crud is the minimum base interface to a data store
type Crud interface {
	Creater
	Reader
	Updater
	Destroyer
}

// Creater is the interface to write values
type Creater interface {
	Create(bucket, key string, target encoding.BinaryMarshaler) error
}

// Reader is the interface to read values
type Reader interface {
	Read(bucket, key string, target encoding.BinaryUnmarshaler) (err error)
}

// Updater is the interface to change values
type Updater interface {
	Update(bucket, key string, target encoding.BinaryMarshaler) error
}

// Destroyer is the interface to delete values
type Destroyer interface {
	Destroy(bucket, key string) error
}

type BinaryIterator interface {
	Next() bool
	Value(value encoding.BinaryUnmarshaler) (key string, err error)
}

// Lister is the interface to list what's in the data store
type Lister interface {
	ListBuckets() ([]string, error)
	ListKeys(bucket string) ([]string, error)
}

type Iterator interface {
	Lister
	Iter(bucket string) BinaryIterator
}
