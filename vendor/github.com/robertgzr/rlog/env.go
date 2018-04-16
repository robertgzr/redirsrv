package rlog

import (
	"os"
	"strings"

	apex "github.com/apex/log"
)

const defaultEnvKey = "GO_LOG"

func ParseEnv() {
	std.Warn("Calling `ParseEnv` explicitly is no longer neccessary")
}

func parseEnv() LvlMap {
	var lvlmap = LvlMap{"*": apex.DebugLevel}

	env := os.Getenv(defaultEnvKey)
	if env == "" {
		return lvlmap
	}

	return lvlsFromString(env, lvlmap)
}

func lvlsFromString(env string, lvlmap LvlMap) LvlMap {
	lvlpairs := strings.Split(env, ",")
	for _, lvlpair := range lvlpairs {
		r := strings.Split(lvlpair, "=")
		if len(r) < 2 {
			lvl := parseLvl(r[0])
			lvlmap["*"] = lvl
			continue
		}
		pkg := r[0]
		lvl := parseLvl(r[1])
		lvlmap[pkg] = lvl
	}
	return lvlmap
}

func parseLvl(str string) apex.Level {
	if lvl, err := apex.ParseLevel(str); err != nil {
		return apex.DebugLevel
	} else {
		return lvl
	}
}
