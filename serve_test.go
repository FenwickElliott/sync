package sync_test

import (
	"crypto/tls"
	"net/http"
	"testing"
)

var client *http.Client

func init() {
	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func TestIn(t *testing.T) {
	var inTests = []struct {
		input      string
		statusCode int
		location   string
	}{
		{"https://cookies.fenwickelliott.io/in?partner=inception&cookie=f9d642125d6d4d660642527b6588b25a88b0d28e", 307, "/forward?rock=8d98eb11968858637053c1bf2f35d8e60dfbc878&glam=fc69f2eb116c7debb15812c107436942e8d99d5c&inception=f9d642125d6d4d660642527b6588b25a88b0d28e&back=https://cookies.fenwickelliott.io"},
	}

	for _, tt := range inTests {
		resp, err := client.Head(tt.input)
		check(err, t)
		if resp.StatusCode != tt.statusCode {
			t.Errorf("/in wrong status, expected: %d, got: %d", tt.statusCode, resp.StatusCode)
		}
		location, err := resp.Location()
		check(err, t)
		if location.RequestURI() != tt.location {
			t.Errorf("/in wrong location, expected: %s, got: %s", tt.location, location.RequestURI())
		}
	}
}

func check(err error, t *testing.T) {
	if err != nil {
		t.Errorf(err.Error())
	}
}
