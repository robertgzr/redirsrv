package bunt

import (
	"encoding"
	"strings"

	"github.com/robertgzr/kiwi"
	bunt "github.com/tidwall/buntdb"
)

func (c *Client) ListKeys(b string) ([]string, error) {
	var keys []string

	err := c.View(func(tx *bunt.Tx) error {
		return tx.AscendKeys(b+":*", func(k, _ string) bool {
			keys = append(keys, strings.TrimPrefix(k, b+":"))
			return true
		})
	})
	return keys, err
}

func (c *Client) ListBuckets() ([]string, error) {
	var buckets = make([]string, 0, len(c.Buckets))
	for b, _ := range c.Buckets {
		buckets = append(buckets, b)
	}
	return buckets, nil
}

type buntIterator struct {
	data [][2]string
	pos  int
}

func (it *buntIterator) Next() bool {
	it.pos += 1
	if it.pos >= len(it.data) {
		return false
	}
	return true
}

func (it *buntIterator) Value(t encoding.BinaryUnmarshaler) (string, error) {
	return it.data[it.pos][0], t.UnmarshalBinary([]byte(it.data[it.pos][1]))
}

func (c *Client) Iter(b string) kiwi.BinaryIterator {
	var data [][2]string

	if err := c.DB.View(func(tx *bunt.Tx) error {
		return tx.AscendKeys(b+":*", func(k, v string) bool {
			data = append(data, [2]string{strings.TrimPrefix(k, b+":"), v})
			return true
		})
	}); err != nil {
		return nil
	}
	return &buntIterator{
		data: data,
		pos:  -1,
	}
}
