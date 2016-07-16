// Copyright (c) 2016 - Sarjono Mukti Aji <me@simukti.net>
// Unless otherwise noted, this source code license is MIT-License

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	args := []string{"-lsah", os.TempDir()}
	cmd := exec.Command("/bin/ls", args...)
	out, outErr := cmd.Output()

	if outErr != nil {
		log.Fatalln(outErr)
	}

	fmt.Println(string(out))
}
