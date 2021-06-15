package gocord

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const baseurl = "https://discordapp.com/api/v9"

func (c *Client) sendRequest(endpoint, method, body string) (string, error) {
	var rBody io.Reader = nil
	if body != "" {
		rBody = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, baseurl+endpoint, rBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("authorization", c.Options.Token)
	if body != "" {
		req.Header.Set("content-type", "application/json")
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
