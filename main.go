package main

import (
	"blockexchange/api"
	"blockexchange/db"
	"blockexchange/jobs"
	"blockexchange/types"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	if os.Getenv("LOGLEVEL") == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Info("Starting")
	cfg := types.CreateConfig()
	err := cfg.Validate()
	if err != nil {
		panic(err)
	}

	repos, err := db.Init()
	if err != nil {
		panic(err)
	}

	// populate database with test data (users, tokens)
	if os.Getenv("BLOCKEXCHANGE_TEST_DATA") == "true" {
		err = db.PopulateTestData(repos)
		if err != nil {
			panic(err)
		}
	}

	// set up server
	api, router, err := api.NewApi(repos, cfg)
	if err != nil {
		panic(err)
	}

	// main entry
	http.Handle("/", handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router))

	// prometheus metrics, on a separate (internal) listener if configured
	if cfg.MetricsAddr != "" {
		metrics_mux := http.NewServeMux()
		metrics_mux.Handle("/metrics", promhttp.Handler())
		go func() {
			logrus.WithField("addr", cfg.MetricsAddr).Info("Serving metrics")
			err := http.ListenAndServe(cfg.MetricsAddr, metrics_mux)
			if err != nil {
				panic(err)
			}
		}()
	} else {
		http.Handle("/metrics", promhttp.Handler())
	}
	server := &http.Server{Addr: ":8080", Handler: nil}

	// start background jobs
	jobs.Start(repos, api)

	go func() {
		// listen to web requests
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	var captureSignal = make(chan os.Signal, 1)
	signal.Notify(captureSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-captureSignal
	logrus.Info("Preparing shutdown")

	//stop api
	api.Stop()
	time.Sleep(5 * time.Second)
	logrus.Info("Shutdown complete")
	server.Shutdown(context.Background())
}
