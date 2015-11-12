package handlers

import (
	"net/http"
	"encoding/json"
	log "github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/BenoitHanotte/shorturls/urlhelper"
)

// the structure of a request
type create_request_body struct {
	Url		string		// the url to shorten
	Perso	string		// the requested personalisation, CAN BE NOT SET
}

// factory to create the handler
func CreateHandler(reachTimeoutMs int) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		// log request for debugging purposes (eg: crash, ...)
		log.WithField("request", r).Debug("create request received")

		// Unmarhall JSON to structure
		decoder := json.NewDecoder(r.Body)

		// unmarshall JSON
		var body create_request_body
		err := decoder.Decode(&body)
		if err != nil {
			log.WithError(err).Error("can not unmarshall JSON body of create request, returning 400: Bad Request")
			// return a 400: Bad Request response
			w.WriteHeader(400)
			return
		}

		// check that the url exists in the structure, and that it is correct
		if body.Url == "" || !urlhelper.IsValid(body.Url) {
			log.Error("incorrect url in body of create request, returning 400: Bad Request")
			w.WriteHeader(400)
			return
		}

		// check that URL is reachable (no intranet, no tor url, ...)
		if !urlhelper.IsReachable(body.Url, reachTimeoutMs) {
			log.WithField("url", body.Url).Error("unreachable URL submitted, returning 400 bad request")
			w.WriteHeader(400)
			return
		}

		w.Write([]byte("create\n"))
	}
}

