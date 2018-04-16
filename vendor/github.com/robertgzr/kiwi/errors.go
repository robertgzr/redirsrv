package kiwi

import (
	"fmt"
)

func IsNotFound(e error) bool {
	_, ok := e.(*notFoundError)
	return ok
}

type notFoundError struct {
	typ  string
	name string
}

func (e *notFoundError) Error() string {
	return fmt.Sprintf("%s %q not found", e.typ, e.name)
}

func NewKeyNotFoundError(key string) error {
	return &notFoundError{
		typ:  "key",
		name: key,
	}
}
func NewBucketNotFoundError(bucket string) error {
	return &notFoundError{
		typ:  "bucket",
		name: bucket,
	}
}
