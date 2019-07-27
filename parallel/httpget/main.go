package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	client     *http.Client
	clientOnce sync.Once
)

func httpClient() *http.Client {
	clientOnce.Do(func() {
		transport, ok := http.DefaultTransport.(*http.Transport)
		if !ok {
			log.Fatal("infra: defaulTransport is not *http.Transport")
		}

		transport.DisableKeepAlives = true
		transport.MaxIdleConns = 254
		transport.MaxIdleConnsPerHost = 254
		transport.MaxConnsPerHost = 254
		transport.ResponseHeaderTimeout = time.Second * 30
		transport.IdleConnTimeout = time.Second * 30
		transport.TLSHandshakeTimeout = time.Second * 30
		transport.DialContext = (&net.Dialer{
			Timeout:   time.Second * 30,
			KeepAlive: time.Second * 30,
			DualStack: true,
		}).DialContext

		client = &http.Client{
			Timeout:   time.Second * 30,
			Transport: transport,
		}
	})

	return client
}

var errTimeout = errors.New("request canceled")

type fetchResult struct {
	url       string      //
	content   interface{} // hasil request
	err       error
	totalTime float64
}

// doFetch do actual http call and will use given ctx as http request context.
func doFetch(ctx context.Context, url string, wg *sync.WaitGroup, result chan<- fetchResult) {
	start := time.Now()
	defer wg.Done()
	var (
		err  error
		req  *http.Request
		res  *http.Response
		body []byte
	)

	resChan := fetchResult{url: url, err: err}
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		resChan.totalTime = time.Since(start).Seconds()
		resChan.err = err
		result <- resChan
		return
	}

	res, err = httpClient().Do(req.WithContext(ctx))
	if err != nil {
		if ctx.Err() != nil {
			err = errTimeout
		}
		resChan.totalTime = time.Since(start).Seconds()
		resChan.err = err
		result <- resChan
		return
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		resChan.totalTime = time.Since(start).Seconds()
		resChan.err = err
		result <- resChan
		return
	}

	var r interface{}
	err = json.Unmarshal(body, &r)
	resChan.totalTime = time.Since(start).Seconds()
	resChan.err = err
	resChan.content = fmt.Sprintf("%v", r)

	result <- resChan
}

func parallelFetch(urls []string /*atau []inputStruct*/, timeout time.Duration) []fetchResult {
	resChan := make(chan fetchResult, len(urls))
	wg := sync.WaitGroup{}

	for _, u := range urls {
		wg.Add(1)
		ctx, _ := context.WithTimeout(context.Background(), timeout)
		go doFetch(ctx, u, &wg, resChan)
	}

	// https://stackoverflow.com/questions/46560204/why-does-my-code-work-correctly-when-i-run-wg-wait-inside-a-goroutine
	go func() {
		wg.Wait()
		// if we not close this, the next range over resChan channel will wait forever.
		close(resChan)
	}()

	var result []fetchResult
	for c := range resChan {
		result = append(result, c)
	}

	return result
}

func main() {
	var urls = []string{
		"https://ip4.seeip.org/json",
		"https://httpbin.org/delay/2",
		"https://api.my-ip.io/ip.json",
	}

	timeout := time.Duration(3 * time.Second)
	log.Println(fmt.Sprintf("fetch start with timeout: %.f second", timeout.Seconds()))

	res := parallelFetch(urls, timeout)
	for _, r := range res {
		log.Println(fmt.Sprintf("\t> %s (done in: %.2f) --> err: %v", r.url, r.totalTime, r.err))
	}

	log.Println("fetch done")
}
