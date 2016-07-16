// Copyright (c) 2016 - Sarjono Mukti Aji <me@simukti.net>
// Unless otherwise noted, this source code license is MIT-License

package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Result struct {
	filename string
	url      string
	err      error
	file     *os.FileInfo
	time     float64
}

var (
	downloadTimeout = time.Duration(3)
)

func main() {
	urls := []string{
		"https://httpbin.org/robots.txt",
		"https://httpbin.org/image/png",
		"https://httpbin.org/image/jpeg",
	}

	download(urls)
}

func download(urls []string) []*Result {
	recv := make(chan *Result)
	results := []*Result{}

	for _, u := range urls {
		go func(u string) {
			start := time.Now()
			url, _ := url.Parse(u)
			basename := filepath.Base(url.Path)
			res := &Result{url: u, filename: basename}
			fname := filepath.Join(os.TempDir(), basename)
			o, oErr := os.Create(fname)

			if oErr != nil {
				res.err = oErr
			}
			defer o.Close()

			done := make(chan error, 1)

			go func() {
				fmt.Println(fmt.Sprintf("Downloading %s, started at %s", u, time.Now().String()))
				r, rErr := http.Get(u)
				if rErr != nil {
					os.Remove(fname)
					done <- rErr
				}

				if r.StatusCode != http.StatusOK {
					os.Remove(fname)
					done <- errors.New(fmt.Sprintf("%s inaccessible", u))
				}
				defer r.Body.Close()

				if _, dErr := io.Copy(o, r.Body); dErr != nil {
					os.Remove(fname)
					done <- dErr
				}

				close(done)
			}()
			select {
			case err := <-done:
				if err != nil {
					res.err = err
				}
			case <-time.After(downloadTimeout * time.Second):
				res.err = errors.New(fmt.Sprintf("\t>>> Download timeout : %s", u))
			}

			f, _ := os.Stat(fname)
			res.file = &f
			end := time.Now()
			res.time = end.Sub(start).Seconds()

			if res.err != nil {
				fmt.Println(res.err.Error())
			} else {
				fmt.Println(fmt.Sprintf("%s downloaded in %f seconds", res.filename, res.time))
			}

			recv <- res
		}(u)
	}

	// loop until all url have a result or timeout
	for {
		select {
		case res := <-recv:
			results = append(results, res)
			if len(results) == len(urls) {
				return results
			}
		}
	}

	return results
}
