package main

import (
	"net/http"

	"github.com/apex/log"
	"github.com/go-chi/chi"
	"github.com/robertgzr/kiwi"
)

const redirKey = "redirkey"

func redirHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, redirKey)
	if key == "" {
		log.Warn("no key found")
		http.Error(w, txt(internalerror, cfg), http.StatusInternalServerError)
		return
	}

	log.WithField("key", key).Debug("looking up")

	var redir kiwi.StringValue
	if err := db.Read("redirs", key, &redir); err != nil {
		log.WithError(err).Error("storage error")

		if kiwi.IsNotFound(err) {
			notFoundHandler(w, r)
			return
		}
		http.Error(w, txt(internalerror, cfg), http.StatusInternalServerError)
		return
	}

	log.WithField("url", redir).Info("success")

	http.Redirect(w, r, string(redir), http.StatusFound)
}
