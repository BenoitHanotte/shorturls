package handlers

import (
	"encoding/json"
	log "github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/gopkg.in/redis.v3"
	"github.com/BenoitHanotte/shorturls/confighelper"
	"net/http"
)

// the structure of a response
type admin_response_body struct {
	Url          string `json:"url"`
	CreationTime string `json:"creationTime"`
	Count        string `json:"count"`
}

// factory to create the handler
func AdminHandler(redisClient *redis.Client, conf *confighelper.Config) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// debug log
		log.WithField("request", r).Debug("admin request received")

		// get the path variable to get the token
		vars := mux.Vars(r)
		token := vars["token"]

		// get the redirection url for this token
		value, err := redisClient.HGetAllMap(token).Result()
		if err != nil && err.Error() != "redis: nil" {
			log.WithError(err).Error("error while retrieving the token infos from redis")
			w.WriteHeader(500) // server error
			return
		} else if value == nil {
			// "redis: nil" is the error is the key is not found
			log.WithField("token", token).Info("token not found")
			w.WriteHeader(404) // not found
			return
		}

		log.WithFields(log.Fields{
			"token": token,
			"value": value}).Debug("mapped values retrieved")

		response := admin_response_body{
			Url:          value["url"],
			CreationTime: value["creationTime"],
			Count:        value["count"],
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(response)

		// avoid caching the page on the client side to not bias the counts
		w.Header().Set("cache-control", "private, max-age=0, no-cache")

		log.WithField("token", token).Info("admin request served")
	}
}
