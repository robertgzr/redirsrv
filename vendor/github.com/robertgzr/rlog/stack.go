package rlog

import (
	"path/filepath"
	"runtime"
)

func getPackage() string {
	var (
		c      = ""
		ok     = false
		offset = 6
	)
	for !ok {
		c = getCaller(offset)
		switch c {
		case "runtime":
			offset--
		case "rlog", "log":
			offset++
		default:
			ok = true
		}
	}
	return c
}

func getCaller(offset int) string {
	_, f, _, ok := runtime.Caller(offset)
	if !ok {
		return ""
	}
	return filepath.Base(filepath.Dir(f))
}
