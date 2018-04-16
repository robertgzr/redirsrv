package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/apex/log"
	"github.com/robertgzr/kiwi"
	"github.com/stevenroose/gonfig"
)

var cfg Config

type Config struct {
	LegacyLinkfilePath string                 `id:"legacy-linkfile"`
	Routes             map[string]interface{} `id:",nohelp"`

	Host string `desc:"host to bind to" default:"localhost"`
	Port int    `desc:"port to bind to" default:"8000"`
}

func (cfg *Config) init() {
	if err := gonfig.Load(cfg, gonfig.Conf{
		FileDefaultFilename: "linkfile.toml",
		FileDecoder:         gonfig.DecoderTOML,
	}); err != nil {
		log.WithError(err).Warn("failed to load linkfile.toml")
	}

	if cfg.LegacyLinkfilePath != "" {
		cfg.loadLegacyLinkfile()
	}

	for key, route := range cfg.Routes {
		if r, ok := route.(string); ok {
			var v = kiwi.StringValue(r)
			if err := db.Create("redirs", key, &v); err != nil {
				log.WithError(err).Warn("failed to write route to db")
			}
		}
	}

	log.Debug("config initialized")
}

type legacyLinkfile struct {
	Short string `json:"short"`
	To    string `json:"to"`
}

func (cfg *Config) loadLegacyLinkfile() {
	data, err := ioutil.ReadFile(cfg.LegacyLinkfilePath)
	if err != nil {
		log.WithError(err).Warn("error reading legacy linkfile")
		return
	}

	var lf []legacyLinkfile
	if err := json.Unmarshal(data, &lf); err != nil {
		log.WithError(err).Warn("error decoding legacy linkfile")
		return
	}

	for _, route := range lf {
		cfg.Routes[route.Short] = route.To
	}
}
