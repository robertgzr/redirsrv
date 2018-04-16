package bunt

import (
	"encoding"

	"github.com/robertgzr/kiwi"
	bunt "github.com/tidwall/buntdb"
)

type Client struct {
	*bunt.DB
	Buckets map[string]int
}

func NewClient() *Client {
	db, _ := bunt.Open(":memory:")
	return &Client{
		DB:      db,
		Buckets: make(map[string]int),
	}
}

func (c *Client) Close() error {
	return c.DB.Close()
}

func (c *Client) Create(b, k string, t encoding.BinaryMarshaler) error {
	if err := c.DB.CreateIndex(b, b+":*", bunt.IndexString); err != nil {
		if err != bunt.ErrIndexExists {
			return err
		}
	}

	value, err := t.MarshalBinary()
	if err != nil {
		return err
	}

	if err := c.DB.Update(func(tx *bunt.Tx) error {
		_, _, err := tx.Set(b+":"+k, string(value), nil)
		return err
	}); err != nil {
		if err == bunt.ErrNotFound {
			return kiwi.NewKeyNotFoundError(k)
		}
		return err
	}
	c.Buckets[b]++
	return nil
}

func (c *Client) Read(b, k string, t encoding.BinaryUnmarshaler) error {
	if _, ok := c.Buckets[b]; !ok {
		return kiwi.NewBucketNotFoundError(b)
	}

	return c.DB.View(func(tx *bunt.Tx) error {
		value, err := tx.Get(b + ":" + k)
		if err != nil {
			if err == bunt.ErrNotFound {
				return kiwi.NewKeyNotFoundError(k)
			}
			return err
		}
		return t.UnmarshalBinary([]byte(value))
	})
}

func (c *Client) Update(b, k string, t encoding.BinaryMarshaler) error {
	if _, ok := c.Buckets[b]; !ok {
		return kiwi.NewBucketNotFoundError(b)
	}
	return c.Create(b, k, t)
}

func (c *Client) Destroy(b, k string) error {
	if _, ok := c.Buckets[b]; !ok {
		return kiwi.NewBucketNotFoundError(b)
	}

	if err := c.DB.Update(func(tx *bunt.Tx) error {
		_, err := tx.Delete(b + ":" + k)
		return err
	}); err != nil {
		if err == bunt.ErrNotFound {
			return kiwi.NewKeyNotFoundError(k)
		}
		return err
	}
	c.Buckets[b]--
	return nil
}
