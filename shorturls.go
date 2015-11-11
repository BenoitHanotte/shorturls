package main

import (
	log "github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/BenoitHanotte/shorturls/config"
	"os"
	"fmt"
)

func init() {
	// configure the logger
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
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
