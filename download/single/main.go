// Copyright (c) 2016 - Sarjono Mukti Aji <me@simukti.net>
// Unless otherwise noted, this source code license is MIT-License

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	u := "https://httpbin.org/robots.txt"
	url, _ := url.Parse(u)
	fname := filepath.Join(os.TempDir(), filepath.Base(url.Path))
	o, oErr := os.Create(fname)
	if oErr != nil {
		log.Fatalln(oErr)
	}
	defer o.Close()

	r, rErr := http.Get(u)
	if rErr != nil {
		log.Fatalln(rErr)
		os.Remove(fname)
	}

	if r.StatusCode != http.StatusOK {
		log.Fatalln("Remote file inaccessible")
		os.Remove(fname)
	}
	defer r.Body.Close()

	if _, dErr := io.Copy(o, r.Body); dErr != nil {
		log.Fatalln(dErr)
		os.Remove(fname)
	}

	f, _ := os.Stat(fname)

	fmt.Println(fmt.Sprintf("Filename\t: %s\nSize\t\t: %d byte\nFiletime\t: %s", f.Name(), f.Size(), f.ModTime().String()))
	os.Remove(fname)
}
