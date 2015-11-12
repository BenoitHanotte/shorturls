package urlhelper

import (
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/asaskevich/govalidator"
	"net/http"
	"regexp"
	"time"
)

func IsValid(url string) bool {
	matched := govalidator.IsURL(url)
	prefixMatch, _ := regexp.MatchString("^https?://", url)
	return matched && prefixMatch
}

// check if URL is reachable on the internet
func IsReachable(url string, reachTimeoutMs int) bool {
	client := http.Client{
		Timeout: time.Duration(reachTimeoutMs) * time.Millisecond}
	_, err := client.Head(url)
	return err == nil
}
