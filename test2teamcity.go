package test2teamcity

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync/atomic"
	"time"
)

const (
	testStarted  = "##teamcity[testStarted name='%s' flowId='%d' captureStandardOutput='false']\n"
	testFinished = "##teamcity[testFinished name='%s' flowId='%d' duration='%d']\n"
	testFailed   = "##teamcity[testFailed name='%s' flowId='%d' message='%s']\n"
	testIgnored  = "##teamcity[testIgnored name='%s' flowId='%d' message='skipped']\n"
	testStdOut   = "##teamcity[testStdOut name='%s' flowId='%d' out='%s']\n"
	testStdErr   = "##teamcity[testStdErr name='%s' flowId='%d' out='%s']\n"
)

func escape(s string) string {
	s = strings.TrimSuffix(s, "\n")
	s = strings.Replace(s, "|", "||", -1)
	s = strings.Replace(s, "\n", "|n", -1)
	s = strings.Replace(s, "\r", "|n", -1)
	s = strings.Replace(s, "'", "|'", -1)
	s = strings.Replace(s, "]", "|]", -1)
	s = strings.Replace(s, "[", "|[", -1)
	return s
}

/*
=== RUN   Test_Secret
=== RUN   Test_Secret/Child
    rule_test.go:19: sss
=== RUN   Test_Secret/Chil2d
--- FAIL: Test_Secret (0.00s)
    --- FAIL: Test_Secret/Child (0.00s)
    --- SKIP: Test_Secret/Child2 (0.00s)
=== RUN   Test_TaxId
--- PASS: Test_TaxId (0.00s)
=== RUN   Test_Phone
--- PASS: Test_Phone (0.00s)
=== RUN   Test_Email
--- PASS: Test_Email (0.00s)
*/

const (
	run   = `=== RUN `
	fail  = `--- FAIL:`
	skip  = `--- SKIP:`
	pass  = `--- PASS:`
	cont  = `=== CONT`
	pause = `=== PAUSE`
)

func Pipe(in io.Reader, stdout io.Writer) error {
	var (
		reader        = bufio.NewReader(in)
		name   string = ""
		flow   uint64 = 0
		seq    uint64 = 0
		cache         = map[string]uint64{}
	)

	uid := func(name string) uint64 {
		if _, ok := cache[name]; !ok {
			cache[name] = atomic.AddUint64(&seq, 1)
		}
		return cache[name]
	}

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		line = bytes.TrimSpace(line)
		max := len(line)

		value := string(line)

		if strings.HasPrefix(value, run) {
			name = strings.TrimSpace(value[len(run):])
			flow = uid(name)
			if _, err = fmt.Fprintf(stdout, testStarted, name, flow); err != nil {
				return err
			}
			continue
		}

		if strings.HasPrefix(value, cont) {
			name = strings.TrimSpace(value[len(cont):])
			flow = uid(name)
			continue
		}

		var (
			isSkip = strings.HasPrefix(value, skip)
			isFail = strings.HasPrefix(value, fail)
			isPass = strings.HasPrefix(value, pass)
		)

		if isPass || isFail || isSkip {
			var (
				sc     = strings.IndexByte(value, ':') + 1
				e      = strings.LastIndexByte(value, '(')
				d, err = time.ParseDuration(value[e+1 : max-1])
			)
			if err != nil {
				return err
			}

			name = strings.TrimSpace(value[sc:e])
			flow = uid(name)

			switch {
			case isFail:
				if _, err = fmt.Fprintf(stdout, testFailed, name, flow, ""); err != nil {
					return err
				}
			case isSkip:
				if _, err = fmt.Fprintf(stdout, testIgnored, name, flow); err != nil {
					return err
				}
			}
			if _, err = fmt.Fprintf(stdout, testFinished, name, flow, d.Milliseconds()); err != nil {
				return err
			}
			continue
		}
		if _, err = fmt.Fprintf(stdout, testStdOut, name, flow, escape(value)); err != nil {
			return err
		}
	}
	return nil
}
