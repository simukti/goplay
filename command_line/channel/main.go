// Copyright (c) 2016 - Sarjono Mukti Aji <me@simukti.net>
// Unless otherwise noted, this source code license is MIT-License

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	commandTimeout = time.After(3600 * time.Second)
)

func main() {
	args := []string{"-lash", os.TempDir()}
	cmd := exec.Command("/bin/ls", args...)
	done := make(chan error, 1)
	result := make(chan []byte, 1)

	go func() {
		out, err := cmd.Output()
		done <- err
		result <- out
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Fatalln(err)
		}
	case <-commandTimeout:
		if kErr := cmd.Process.Kill(); kErr != nil {
			log.Fatalln(kErr)
		}

		log.Fatalln("Upss...")
	}

	buf := bytes.NewBuffer(<-result)
	fmt.Println(buf.String())
}
