package gocord

import (
	"github.com/valyala/fasthttp"
)

const baseurl = "https://discordapp.com/api/v9"

func (c *Client) sendRequest(endpoint, method, body string) (string, error) {
	res := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(res)

	res.SetRequestURI(baseurl+endpoint)
	
	res.Header.SetMethod(method)
	res.AppendBodyString(body)

	res.Header.Set("authorization", c.Options.Token)
	if body != "" {
		res.Header.SetContentType("application/json")
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := fasthttp.Do(res, resp)
	if err != nil {
		return "", err
	}

	return string(resp.Body()), nil
}
