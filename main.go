package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/pressly/chi"

	"github.com/robertgzr/kiwi"
	"github.com/robertgzr/kiwi/bunt"
	"github.com/robertgzr/rlog"
)

const (
	index = `
    WHAT IS THIS?
    a really simple redirection service

    WHAT CAN IT DO?
    if your're authorized, query /adm?bucket=redirs for all available routes
	`
	forbidden = `
    UNAUTHORIZED ACCESS!
    request assistance / access at r@gnzler.io
	`
	notfound = `
    ROUTE NOT FOUND!
    seems you've tried to access an invalid route

    BUT THIS SHOULD WORK!?
    if you think this is a mistake, contact me at r@gnzler.io
	`
	internalerror = `
    INTERNAL ERROR :(
    oops this was not supposed to happen, tell me about it at r@gnzler.io
	`
)

var (
	host = "localhost"
	port = 8080

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
	r.Get("/", indexHandler)
	r.Get("/:"+redirKey, redirHandler)

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
	w.Write([]byte(index))
}
