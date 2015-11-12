package config

import (
	"errors"
	log "github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/spf13/viper"
	"os"
)

type Config struct {
	ValueLength	int		// the length of the value (eg: x8f9Rz for toto.com/x8f9Rz)
	Host		string	// the host to use (eg: toto.com), default: HOST env variable
	Port 		int		// the port of the server
	Proto		string	// the protocol
	RedisAddr 	string 	// the IP address of the redis node
	RedisPort 	int    	// the port of the redis node
	RedisProto	string	// the protocol for communication with redis
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
		log.WithField("filename", filename).Error("config file not found")
		return nil, errors.New("file not found")
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

	config := Config{
		ValueLength:	viper.GetInt("valueLength"),
		Host:			viper.GetString("host"),
		Port:			viper.GetInt("port"),
		Proto:			viper.GetString("proto"),
		RedisAddr: 		viper.GetString("redisHost"),
		RedisPort: 		viper.GetInt("redisPort"),
		RedisProto: 	viper.GetString("redisProto")}

	log.WithField("config", config).Debug("configuration loaded from YAML")

	return &config, nil
}
