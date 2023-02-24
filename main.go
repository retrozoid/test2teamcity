package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/retrozoid/test2teamcity"
)

// === RUN   Test_Secret
// === RUN   Test_Secret/Child
//     rule_test.go:19: sss
// --- FAIL: Test_Secret (0.00s)
//     --- FAIL: Test_Secret/Child (0.00s)
// === RUN   Test_TaxId
// --- PASS: Test_TaxId (0.00s)
// === RUN   Test_Phone
// --- PASS: Test_Phone (0.00s)
// === RUN   Test_Email
// --- PASS: Test_Email (0.00s)
// FAIL

func main() {
	var (
		args   = flag.Args()
		reader io.ReadCloser
		err    error
	)
	flag.Parse()
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
