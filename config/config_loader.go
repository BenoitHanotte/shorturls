package config

import (
	"errors"
	log "github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/spf13/viper"
)

type Config struct {
	RedisAddr string // the IP address of the redis node
	RedisPort int    // the port of the redis node
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
		return nil, errors.New("config file not found")
	}

	// set defaults
	viper.SetDefault("RedisPort", 1234)

	c := Config{
		RedisAddr: viper.GetString("essai"),
		RedisPort: viper.GetInt("RedisPort")}

	log.WithFields(log.Fields{
		"RedisAddr": c.RedisAddr,
		"RedisPort": c.RedisPort}).Debug("configuration loaded from YAML")

	return &c, nil
}
