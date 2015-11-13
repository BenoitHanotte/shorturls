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

		// generate the token
		token, err := generateFreeToken(body.Token, conf.TokenLength, checkTokenExists(redisClient))
		if err != nil {
			log.WithError(err).Error("can not create a token, aborting")
			w.WriteHeader(500)
			return
		}

		// save to redis with the token as the key
		_, err = redisClient.HMSet(token,
			"url", body.Url,
			"creationTime", strconv.FormatInt(time.Now().Unix(), 10),
			"count", "0").Result()
		if err != nil {
			log.WithError(err).Error("can not set the entry in Redis, aborting")
			w.WriteHeader(500)
			return
		}

		// set expiration time in
		redisClient.ExpireAt(token, time.Now().AddDate(0, conf.ExpirationTimeMonths, 0))

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
func validateToken(token string, tokenLength int) {
	match, _ := regexp.MatchString("^[0-9a-zA-Z]{0,"+strconv.Itoa(tokenLength)+"}$", token)
	return match
}

// generate a token
// CAUTION: if the token is already present, generate a new one
// If more than 3 successive collisions, use one more random character
func generateFreeToken(
	suggestion string,
	tokenLength int,
	checkTokenExistsFunc func (token string) (bool, error)) (string, error) {

	// offset is the number of random characters generated at the end of the suggestion
	var offset = mathhelper.Max(0, tokenLength-len(suggestion))
	var token = ""

	for i := 0; i < maxRetries; i++ {
		token = makeRandString(suggestion, tokenLength, offset)

		// check redis to see if this token already exists -> in that case generate a new one
		exists, err := checkTokenExistsFunc(token)
		if err != nil {
			return "", err
		}

		if !exists {
			// exit loop since there is no collision this time
			return token, nil
		}

		log.WithFields(log.Fields{
			"token":  token,
			"retry":  i,
			"offset": offset}).Debug("collision while generating new token")

		// if already exists, generate another token and try again until a correct token is generated
		if (i == 0 || i%retriesToRaiseOffset == 0) && offset < tokenLength {
			offset += 1
		}
	}
	return "", errors.New("maximum number of retries reached")
}

// factory to create a function checking if a token is already in Redis
// this is necessary to mock this function for the tests
func checkTokenExists(redisClient *redis.Client) func (token string) (bool, error) {
	return func (token string) (bool, error) {
		return redisClient.Exists(token).Result()
	}
}

// make a random string from a starting string and an offset
// (the number of random char at the end)
func makeRandString(start string, length int, offset int) string {
	var token string
	// generate a token
	if len(start) == 0 {
		token = randStringBytesRmndr(offset)
	} else {
		token = start[:length -offset] + randStringBytesRmndr(offset)
	}
	return token
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
