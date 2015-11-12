package urlhelper

import (
	"testing"
)

var validUrls = []string{
	"http://foo.com/blah_blah",
	"http://foo.com/blah_blah/",
	"http://foo.com/blah_blah_(wikipedia)",
	"http://foo.com/blah_blah_(wikipedia)_(again)",
	"http://www.example.com/wpstyle/?p=364",
	"https://www.example.com/foo/?bar=baz&inga=42&quux",
	"http://userid:password@example.com:8080",
	"http://userid:password@example.com:8080/",
	"http://userid@example.com",
	"http://userid@example.com/",
	"http://userid@example.com:8080",
	"http://userid@example.com:8080/",
	"http://userid:password@example.com",
	"http://userid:password@example.com/",
	"http://142.42.1.1/",
	"http://142.42.1.1:8080/",
	"http://foo.com/blah_(wikipedia)#cite-1",
	"http://foo.com/blah_(wikipedia)_blah#cite-1",
	"http://foo.com/(something)?after=parens",
	"http://code.google.com/events/#&product=browser",
	"http://j.mp",
	"http://foo.bar/?q=Test%20URL-encoded%20stuff",
	"http://1337.net",
	"http://a.b-c.de",
	"http://223.255.255.254",
}

var invalidUrls = []string {
	"http://",
	"http://.",
	"http://..",
	"http://../",
	"http://?",
	"http://??",
	"http://??/",
	"http://#",
	"http://##",
	"http://##/",
	"http://foo.bar?q=Spaces should be encoded",
	"//",
	"//a",
	"///a",
	"///",
	"http:///a",
	"foo.com",
	"rdar://1234",
	"h://test",
	"http:// shouldfail.com",
	":// should fail",
	"http://foo.bar/foo(bar)baz quux",
	"ftps://foo.bar/",
	"http://-error-.invalid/",
	"http://a.b--c.de/",
	"http://-a.b.co",
	"http://a.b-.co",
	"http://0.0.0.0",
	"http://.www.foo.bar/",
	"http://www.foo.bar./",
	"http://.www.foo.bar./",
}

func TestIsValid(t *testing.T)  {
	for _, url := range validUrls {
		matches := IsValid(url)
		if !matches {
			t.Error("For", url)
		}
	}

	for _, url := range invalidUrls {
		matches := IsValid(url)
		if matches {
			t.Error("For", url)
		}
	}
}

