// ADDED BY DROP - https://github.com/matryer/drop (v0.7)
//  source: bitbucket.org/robertgzr/drops /ulid (628d7973d726ea3d41438c387bc358750600e45f)
//  update: drop -f bitbucket.org/robertgzr/drops ulid
// license:  (see repo for details)

// Pakcgae ulid contains a very simple helper to generate ULIDs
package main

import (
	"math/rand"

	"github.com/oklog/ulid"
)

func NowULID() ulid.ULID {
	ts := ulid.Now()
	random := rand.New(rand.NewSource(int64(ts)))
	return ulid.MustNew(ts, random)
}
