package main

import (
	log "github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/BenoitHanotte/shorturls/config"
	"os"
	"fmt"
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
	c, err := config.LoadConfigYAML("config")
	if err != nil { // no config could be read (eg: bad filename, missing value...)
		log.Error("no correct config found, exit the program")
		return
	}
	log.Info("configuration loaded")

	fmt.Println(c.RedisPort)
}
