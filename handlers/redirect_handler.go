package handlers

import (
	log "github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/gopkg.in/redis.v3"
	"github.com/BenoitHanotte/shorturls/config"
	"net/http"
)

// factory to create the handler
func RedirectHandler(redisClient *redis.Client, conf *config.Config) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		// debug log
		log.WithField("request", r).Debug("redirect request received")

		// get the path variable to get the token
		vars := mux.Vars(r)
		token := vars["token"]

		// get the redirection url for this token
		url, err := redisClient.HGet(token, "url").Result()
		if err != nil && err.Error() != "redis: nil" {
			log.WithError(err).Error("error while retrieving the redirection url from redis")
			w.WriteHeader(500) // server error
			return
		} else if url == "" {
			// "redis: nil" is the error is the key is not found
			log.WithField("token", token).Info("token not found")
			w.WriteHeader(404) // not found
			return
		}
		// consider that url in Redis is correct from here

		// increment count
		count, err := redisClient.HIncrBy(token, "count", 1).Result()
		if err != nil {
			log.WithError(err).Error("error while incrementing count")
			// no server error, we can still redirect the user
		}

		// debug log
		log.WithFields(log.Fields{
			"token": token,
			"url":   url,
			"count": count}).Debug("redirection data retrieved")

		// redirect
		w.Header().Set("Location", url)
		// avoid caching the page on the client side to not bias the counts
		w.Header().Set("cache-control", "private, max-age=0, no-cache")
		w.WriteHeader(301) // moved permanently

		log.WithField("token", token).Info("redirect request served")
	}
}
