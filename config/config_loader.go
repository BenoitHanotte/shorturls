package config

import (
	"errors"
	log "github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/spf13/viper"
	"os"
)

type Config struct {
	ValueLength    	int    // the length of the value (eg: x8f9Rz for toto.com/x8f9Rz)
	ReachTimeoutMs 	int    // the timeout in ms when checking the reachability of an url
	Host           	string // the host to use (eg: toto.com), default: HOST env variable
	Port           	int    // the port of the server
	Proto          	string // the protocol
	RedisAddr      	string // the IP address of the redis node
	RedisPort      	int    // the port of the redis node
	RedisProto     	string // the protocol for communication with redis
	RedisDB   		string // the redis database
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
	if os.Getenv("REDIS_PORT_6379_TCP_PROTO")!="" {
		viper.Set("redisProto", os.Getenv("REDIS_PORT_6379_TCP_PROTO"))
	}
	if os.Getenv("REDIS_DB")!="" {
		viper.Set("redisDB", os.Getenv("REDIS_DB"))
	}

	config := Config{
		ValueLength:	viper.GetInt("valueLength"),
		ReachTimeoutMs:	viper.GetInt("reachTimeoutMs"),
		Host:			viper.GetString("host"),
		Port:			viper.GetInt("port"),
		Proto:			viper.GetString("proto"),
		RedisAddr: 		viper.GetString("redisHost"),
		RedisPort: 		viper.GetInt("redisPort"),
		RedisProto: 	viper.GetString("redisProto"),
		RedisDB:	 	viper.GetString("redisDB")}

	log.WithField("config", config).Debug("configuration loaded from YAML")

	return &config, nil
}
