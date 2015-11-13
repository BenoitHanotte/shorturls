package handlers

import (
	"encoding/json"
	log "github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/gopkg.in/redis.v3"
	"github.com/BenoitHanotte/shorturls/config"
	"github.com/BenoitHanotte/shorturls/mathhelper"
	"github.com/BenoitHanotte/shorturls/urlhelper"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"errors"
)

// constant for the token generation
const retriesToRaiseOffset = 3 // if after X retries there is still token collision,
// make random part of the token longer by 1 character
const maxRetries = 20 // the number of retries before giving up (too many collisions)

// available characters to generate random strings
const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// the structure of a request (unmarshalled from JSON)
type create_request_body struct {
	Url   string // the url to shorten
	Token string // the requested personalisation, CAN BE NOT SET
}

// the structure of a response (marshalled to JSON)
type create_response_body struct {
	Url string `json:"url"` // the url, marshalled to "url" and not "Url"
}

func init() {
	// seed the random number generator
	rand.Seed(time.Now().UnixNano())
}

// factory to create the handler
func CreateHandler(redisClient *redis.Client, conf *config.Config) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		// log request for debugging purposes (eg: crash, ...)
		log.WithField("request", r).Debug("create request received")

		// Unmarshall JSON to structure
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

		log.WithField("request", r).Debug("create request received")

		// check that the url exists in the structure, and that it is correct
		if body.Url == "" || !urlhelper.IsValid(body.Url) {
			log.Error("incorrect url in body of create request, returning 400: Bad Request")
			w.WriteHeader(400)
			return
		}

		// check that URL is reachable (no intranet, no tor url, ...)
		if !urlhelper.IsReachable(body.Url, conf.ReachTimeoutMs) {
			log.WithField("url", body.Url).Error("unreachable URL submitted, returning 400 bad request")
			w.WriteHeader(400)
			return
		}

		// validate the suggestion (only letters and digits)
		if !validateToken(body.Token, conf.TokenLength) {
			log.WithField("token", body.Token).Error("invalid custom token, aborting")
			w.WriteHeader(400)
			return
		}

		// create a random token generator
		randomTokenGenerator := randomTokenGenerator(body.Token, conf.TokenLength);
		var token string

		// try to insert with a new random token as long as the lock on the token cannot be acquired
		for  {
			token, err = randomTokenGenerator()
			if err != nil {
				// we tried to generate too many token, abort
				log.WithError(err).Error("too many collisions for token, aborting")
				w.WriteHeader(500)
				return
			}

			// use HSetNX to get lock on the Token
			lockAcquired, err := redisClient.HSetNX(token, "url", body.Url).Result()
			if lockAcquired {
				// debug log
				log.WithField("token", token).Debug("lock was acquired")

				// lock could be acquired: we reserved the token !
				// proceed by setting other fields
				_, err = redisClient.HMSet(token, "creationTime", strconv.FormatInt(time.Now().Unix(), 10),
					"count", "0").Result()
				// set expiration time in 3 months
				redisClient.ExpireAt(token, time.Now().AddDate(0, conf.ExpirationTimeMonths, 0))
				break;	// leave the loop: don't try to generate a new token
			}

			// if there was an error while setting in the map (more than just key already present)
			if err != nil {
				log.WithError(err).Error("can not set the entry in Redis, aborting")
				w.WriteHeader(500)
				return
			}

			// debug log
			log.WithField("token", token).Debug("could not acquire lock, retrying if allowed")
		}

		// log success
		log.WithFields(log.Fields{
			"url":   body.Url,
			"token": token}).Info("new short link created")

		// generate response
		response := create_response_body{
			Url: urlhelper.Build(conf.Proto, conf.Host, conf.Port, token),
		}

		w.WriteHeader(201) // return 201: created
		encoder := json.NewEncoder(w)
		encoder.Encode(response)
	}
}

// function to validate the token suggested by the user
func validateToken(token string, tokenLength int) bool {
	match, _ := regexp.MatchString("^[0-9a-zA-Z]{0,"+strconv.Itoa(tokenLength)+"}$", token)
	return match
}

// Factory to create a function which generates random token
func randomTokenGenerator(suggestion string, tokenLength int) func() (string, error) {

	offset := mathhelper.Max(0, tokenLength-len(suggestion))		// nb of char to randomize at the end of the token
	retry := 0														// current retry

	return func() (string, error) {

		if retry >maxRetries {
			return "", errors.New("maximum number of retry reached to generate token")
		}

		// make a random token
		var token string

		if len(suggestion) == 0 {
			token = randStringBytesRmndr(offset)
		} else {
			token = suggestion[:tokenLength - offset] + randStringBytesRmndr(offset)
		}

		// raise offset is too many retries
		if (retry == 0 && offset==0 || retry %retriesToRaiseOffset == 0) && offset < tokenLength {
			offset += 1
		}

		retry++;
		return token, nil
	}

}

// generate random strings of size n
// from: http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func randStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
