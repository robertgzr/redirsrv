package kiwi

import (
	"fmt"
	"strings"
)

func ToString(c Client) string {
	var buf strings.Builder

	buckets, err := c.ListBuckets()
	if err != nil {
		panic(err)
	}
	for _, b := range buckets {
		buf.WriteString(fmt.Sprintf("bucket %q:\n", b))

		var lenv = 0
		for it := c.Iter(b); it.Next(); {
			var v StringValue
			k, err := it.Value(&v)
			if err != nil {
				panic(err)
			}
			buf.WriteString(fmt.Sprintf("  > %s => %v\n", k, v))
			lenv++
		}
		buf.WriteString(fmt.Sprintf("%d entries\n\n", lenv))
	}

	return buf.String()
}
