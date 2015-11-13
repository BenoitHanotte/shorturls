package config

import (
	"errors"
	log "github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/spf13/viper"
	"os"
)

type Config struct {
	TokenLength    			int    			// the length of the value (eg: x8f9Rz for toto.com/x8f9Rz)
	ReachTimeoutMs 			int    			// the timeout in ms when checking the reachability of an url
	ExpirationTimeMonths	int				// the number of months before a short url is deleted
	Host           			string 			// the host to use (eg: toto.com), default: HOST env variable
	Port           			int    			// the port of the server
	Proto          			string 			// the protocol
	RedisHost      			string 			// the host of the redis node
	RedisPort      			int    			// the port of the redis node
	RedisDB        			int 			// the redis database
	RedisPassword  			string 			// the password for redis
}

// Load a YAML config file and put values in a Config object
func LoadConfigYAML(filename string) (*Config, error) {

	log.WithField("filename", filename).Debug("loading configuration from YAML")

	// load config from file from the working directory
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		// log the error
		log.WithFields(log.Fields{
			"filename": filename,
			"error": err}).Error("config file not found")
		return nil, errors.New("config file not found")
	}

	// overrides with environment variables (for debug, or to use docker links)
	if os.Getenv("HOST")!="" {
		viper.Set("host", os.Getenv("HOST"))
	}
	if os.Getenv("PORT")!="" {
		viper.Set("port", os.Getenv("PORT"))
	}
	if os.Getenv("PROTO")!="" {
		viper.Set("proto", os.Getenv("PROTO"))
	}
	if os.Getenv("REDIS_PORT_6379_TCP_PORT")!="" {
		viper.Set("redisPort", os.Getenv("REDIS_PORT_6379_TCP_PORT"))
	}
	if os.Getenv("REDIS_PORT_6379_TCP_ADDR")!="" {
		viper.Set("redisHost", os.Getenv("REDIS_PORT_6379_TCP_ADDR"))
	}
	if os.Getenv("REDIS_DB")!="" {
		viper.Set("redisDB", os.Getenv("REDIS_DB"))
	}
	if os.Getenv("REDIS_PASSWORD")!="" {
		viper.Set("redisPassword", os.Getenv("REDIS_PASSWORD"))
	}

	config := Config{
		TokenLength:	viper.GetInt("tokenLength"),
		ReachTimeoutMs:	viper.GetInt("reachTimeoutMs"),
		Host:			viper.GetString("host"),
		Port:			viper.GetInt("port"),
		Proto:			viper.GetString("proto"),
		RedisHost: 		viper.GetString("redisHost"),
		RedisPort: 		viper.GetInt("redisPort"),
		RedisDB:	 	viper.GetInt("redisDB"),
		RedisPassword:	viper.GetString("redisPassword")}

	return &config, nil
}
