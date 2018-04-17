package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"
	"github.com/go-chi/chi"
	"github.com/oklog/run"
	"github.com/pkg/errors"

	"github.com/robertgzr/kiwi"
	"github.com/robertgzr/kiwi/bunt"
	"github.com/robertgzr/rlog"
)

var (
	db kiwi.Client
)

func init() {
	log.SetHandler(rlog.Default)
	log.SetLevel(log.DebugLevel)
}

func main() {
	db = bunt.NewClient()
	defer db.Close()

	cfg.init()

	r := chi.NewRouter()
	// r.Use(middleware.DefaultLogger)
	r.NotFound(notFoundHandler)

	r.Get("/", indexHandler)
	r.Get("/{"+redirKey+"}", redirHandler)

	bearerToken := NowULID().String()
	log.WithField("bearer", bearerToken).Info("token for API auth")
	r.With(AuthBearer(bearerToken)).Handle("/adm", kiwi.RestHandler(db))

	var g run.Group
	{
		// start the server
		//
		ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
		if err != nil {
			log.WithError(err).Error("error opening listener")
			os.Exit(1)
		}
		g.Add(func() error {
			log.WithFields(log.Fields{
				"host": cfg.Host,
				"port": cfg.Port,
			}).Info("started listening")
			err := http.Serve(ln, r)
			if err != nil {
				return err
			}
			return nil
		}, func(err error) {
			ln.Close()
			log.Warn("stopped listening")
		})
	}
	{
		// signal catcher
		//
		signalC := make(chan os.Signal)
		signal.Notify(signalC, syscall.SIGINT, syscall.SIGHUP)
		g.Add(func() error {
			for {
				s := <-signalC
				log.WithField("signal", s.String()).Debug("received signal")

				switch s {
				case os.Interrupt:
					return errors.Errorf("received signal %v", s)
				case syscall.SIGHUP:
					cfg.init()
				}
			}
		}, func(error) {
			return
		})
	}

	if err := g.Run(); err != nil {
		log.WithField("cause", err.Error()).Warn("reason for shutdown")
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(txt(index, cfg)))
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, txt(notfound, cfg), http.StatusNotFound)
}
