package main

import (
	"net/http"
	"strings"

	"github.com/apex/log"
)

func AuthBearer(token string) func(h http.Handler) http.Handler {
	if token == "" {
		panic("token for Authentification:Bearer can not be empty")
	}

	b := authbearer{tok: token}

	fn := func(h http.Handler) http.Handler {
		b.h = h
		return &b
	}
	return fn
}

type authbearer struct {
	tok string
	h   http.Handler
}

func (a *authbearer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"bucket": r.URL.Query().Get("bucket"),
		"key":    r.URL.Query().Get("key"),
	}).Debug("protected resource")

	auth := r.Header.Get("Authorization")
	if auth == "" {
		log.Warn("unauthorized")
		http.Error(w, txt(forbidden, cfg), http.StatusForbidden)
		return
	}

	idx := strings.Index(auth, ":")
	if idx == -1 {
		log.Warn("unauthorized")
		http.Error(w, txt(forbidden, cfg), http.StatusForbidden)
		return
	}

	bearerToken := auth[idx+1:]
	log.WithField("token", bearerToken).Debug("found bearer token")

	if a.tok != bearerToken {
		log.Warn("unauthorized")
		http.Error(w, txt(forbidden, cfg), http.StatusForbidden)
		return
	}

	a.h.ServeHTTP(w, r)
}
