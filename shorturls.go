package main

import (
	log "github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/BenoitHanotte/shorturls/config"
	"github.com/BenoitHanotte/shorturls/handlers"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
)

func init() {
	// configure the logger
	// The level can be changed with the LOG_LEVEL environment variable, default=info
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	// set the output
	log.SetOutput(os.Stdout)
}

func main() {

	log.Info("starting")

	// Load the configuration from the config.yaml file
	conf, err := config.LoadConfigYAML("config")
	if err != nil { // no config could be read (eg: bad filename, missing value...)
		log.WithField("err", err).Error("incorrect config, exiting")
		return
	}
	log.Info("configuration loaded")


	// create the router
	r := mux.NewRouter()
	// Routes
	var valueRegexp string =  "[0-9a-zA-Z]{"+strconv.Itoa(conf.ValueLength)+"}"

	r.HandleFunc("/{value:"+valueRegexp+"}", handlers.RetrieveHandler).
		Methods("GET").Host(conf.Host)
	r.HandleFunc("/shortlink/{value:"+valueRegexp+"}", handlers.CreateHandler).
		Methods("POST").Headers("Content-Type", "application/json").Host(conf.Host)
	r.HandleFunc("/admin/{value:"+valueRegexp+"}", handlers.AdminHandler).
		Methods("GET").Host(conf.Host)

	// Bind to a port and pass our router in
	log.Info("starting the router...")
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), r)
	if err!=nil {
		log.WithField("err", err).Error("could not start the router, exiting")
		return
	}
}
