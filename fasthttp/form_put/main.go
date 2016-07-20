// Copyright (c) 2016 - Sarjono Mukti Aji <me@simukti.net>
// Unless otherwise noted, this source code license is MIT-License

package main

import (
	"github.com/valyala/fasthttp"
	"log"
	"net/url"
)

func main() {
	form := url.Values{}
	form.Add("key", "value")

	u := "https://httpbin.org/put"
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	req.Header.SetMethodBytes([]byte("PUT"))
	req.Header.SetUserAgentBytes([]byte("X"))
	req.SetRequestURI(u)
	// application/x-www-form-urlencoded"
	req.SetBody([]byte(form.Encode()))

	if httpErr := fasthttp.Do(req, res); httpErr != nil {
		log.Println(httpErr)
	}

	log.Println(res.Header.String())
	log.Println(string(res.Body()))
}
