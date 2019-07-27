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
	urlParse, _ := url.Parse(u)
	fname := filepath.Join(os.TempDir(), filepath.Base(urlParse.Path))
	o, oErr := os.Create(fname)
	if oErr != nil {
		log.Fatalln(oErr)
		return
	}
	defer o.Close()
	defer os.Remove(fname)

	r, rErr := http.Get(u)
	if rErr != nil {
		log.Fatalln(rErr)
		return
	}

	if r.StatusCode != http.StatusOK {
		log.Fatalln("Remote file inaccessible")
		return
	}
	defer r.Body.Close()

	if _, dErr := io.Copy(o, r.Body); dErr != nil {
		log.Fatalln(dErr)
		return
	}

	f, err := os.Stat(fname)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println(fmt.Sprintf("Filename\t: %s\nSize\t\t: %d byte\nFiletime\t: %s", f.Name(), f.Size(), f.ModTime().String()))
}
