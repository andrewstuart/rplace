package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/oauth2"
)

// kinda hacky, but it works
var searchRegex = regexp.MustCompile("\"session\":{\"accessToken\":\"(?P<token>.+)\",\"expires\":\"(?P<expireTime>.+)\",\"expiresIn\":(?P<expiresIn>[0-9]+)")

func getToken(ctx context.Context) (oauth2.Token, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://new.reddit.com/r/place/", nil)
	if err != nil {
		return oauth2.Token{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return oauth2.Token{}, err
	}
	defer res.Body.Close()

	buf := &bytes.Buffer{}
	io.Copy(buf, res.Body)

	if match := searchRegex.FindStringSubmatch(buf.String()); match != nil {
		token := match[searchRegex.SubexpIndex("token")]
		expireTime := match[searchRegex.SubexpIndex("expireTime")]

		exp, _ := time.Parse(time.RFC3339, expireTime)
		return oauth2.Token{
			AccessToken: token,
			Expiry:      exp,
		}, nil
	}

	return oauth2.Token{}, fmt.Errorf("could not find anonymous token")
}
