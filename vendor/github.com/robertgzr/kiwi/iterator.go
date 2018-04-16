package kiwi

import "encoding"

type kiwiIterator struct {
	data interface {
		Reader
		Lister
	}
	bucket string
	keys   []string
	pos    int
}

func GenericIterator(c interface {
	Reader
	Lister
}, bucket string) *kiwiIterator {
	keys, err := c.ListKeys(bucket)
	if err != nil {
		return nil
	}
	return &kiwiIterator{
		data:   c,
		bucket: bucket,
		keys:   keys,
		pos:    -1,
	}
}

func (it *kiwiIterator) Next() bool {
	it.pos += 1
	if it.pos >= len(it.keys) {
		return false
	}
	return true
}

func (it *kiwiIterator) Value(t encoding.BinaryUnmarshaler) (string, error) {
	err := it.data.Read(it.bucket, it.keys[it.pos], t)
	return it.keys[it.pos], err
}
