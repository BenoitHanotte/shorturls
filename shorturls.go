package main

import (
	log "github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/gopkg.in/redis.v3"
	"github.com/BenoitHanotte/shorturls/config"
	"github.com/BenoitHanotte/shorturls/handlers"
	"net/http"
	"os"
	"strconv"
	"fmt"
)

func main() {

	// setUp the logger
	toDefer := setUpLog()
	// the returned function is to defer, used to close log file on exit
	defer toDefer();

	log.Info("starting")

	// Load the configuration from the config.yaml file
	conf, err := config.LoadConfigYAML("config")
	if err != nil {
		log.WithError(err).Fatal("incorrect config, exiting")
		return	// do not exit since the log file still has to be closed by a defered function
	}
	log.Info("configuration loaded")

	// create the redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.RedisHost + ":" + strconv.Itoa(conf.RedisPort),
		Password: conf.RedisPassword,  // no password set
		DB:       int64(conf.RedisDB), // use default DB
	})

	// create the router
	r := mux.NewRouter()
	// Routes
	var valueRegexp string = "[0-9a-zA-Z]{" + strconv.Itoa(conf.TokenLength) + "}"

	r.HandleFunc("/{token:"+valueRegexp+"}", handlers.RedirectHandler(redisClient, conf)).
		Methods("GET").Host(conf.Host)
	r.HandleFunc("/shortlink", handlers.CreateHandler(redisClient, conf)).
		Methods("POST").Headers("Content-Type", "application/json").Host(conf.Host)
	r.HandleFunc("/admin/{token:"+valueRegexp+"}", handlers.AdminHandler(redisClient, conf)).
		Methods("GET").Host(conf.Host)

	// Bind to a port and pass our router in
	log.Info("starting the router...")
	err = http.ListenAndServe(":"+strconv.Itoa(conf.Port), r)
	if err != nil {
		log.WithError(err).Fatal("could not start the router, exiting")
		return
	}
}

// set up the logrus logger
// returns a function to defer, used to defer closing the file in which the logs are written
func setUpLog() func() {
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

	// set the logger output
	logFile := os.Getenv("LOG_FILE")
	switch logFile {
	case "":
		// also when not set (actual default)
		log.SetOutput(os.Stderr)
	default:
		// consider any other value as a filepath
		f, err := os.Create(logFile)
		if err!=nil {
			fmt.Errorf("Unable to open log file: "+err.Error())
		}
		log.SetOutput(f)
		return func() {
			defer f.Close()
		}
	}

	// nothing to defer if log is not written to file
	return func() {}
}