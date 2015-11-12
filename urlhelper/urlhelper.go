package urlhelper

import (
	"github.com/BenoitHanotte/shorturls/Godeps/_workspace/src/github.com/asaskevich/govalidator"
	"net/http"
	"time"
	"strconv"
)

func IsValid(url string) bool {
	if len(url)<10 {
		return false
	}
	matched := govalidator.IsURL(url)
	return matched && (url[:4]=="http" || url[:5]=="https")
}

// check if URL is reachable on the internet
func IsReachable(url string, reachTimeoutMs int) bool {
	client := http.Client{
		Timeout: time.Duration(reachTimeoutMs) * time.Millisecond}
	_, err := client.Head(url)
	return err == nil
}

func Build(proto string, host string, port int, ext string) string {
	if (port==80) {
		return proto+"://"+host+"/"+ext
	} else {
		return proto+"://"+host+":"+strconv.Itoa(port)+"/"+ext
	}
}
