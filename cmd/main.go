package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/retrozoid/test2teamcity"
)

func main() {
	flag.Parse()
	var (
		args   = flag.Args()
		reader io.ReadCloser
		err    error
	)
	if flag.NArg() == 0 {
		reader = os.Stdout
	} else {
		cmd := exec.Command(args[0], args[1:]...)
		reader, err = cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err = cmd.Start(); err != nil {
			log.Fatal(err)
		}
		defer cmd.Wait()
	}
	if err = test2teamcity.Pipe(reader, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
