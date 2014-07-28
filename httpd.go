// httpPerfTest is a simple test server that will serve back arbitrary sized responses
// max speed I have seen so far is ~ 450k req/second using -m 32
// wrk -c 500 -d 10 -t 10
// against localhost
// pprof is enabled and can be reached at http://localhost:8081/debug/pprof/profiling

package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
)

const alphaNum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

type controller struct {
	responseData  []byte
	contentLength string
}

func newController(size int) *controller {
	c := &controller{}
	c.responseData = make([]byte, size)
	for i := 0; i < size; i++ {
		c.responseData[i] = alphaNum[i%len(alphaNum)]
	}
	c.responseData[size-1] = byte('\n')
	c.contentLength = strconv.Itoa(size)
	return c
}

func (c *controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", c.contentLength)
	w.Write(c.responseData)
}

func main() {
	size := flag.Int("s", 128, "response size in bytes")
	httpPort := flag.String("p", "8080", "HTTP listen port")
	profilePort := flag.String("r", "8081", "HTTP pprof listen port")
	maxProcs := flag.Int("m", 1, "GOMAXPROCS")
	flag.Parse()

	runtime.GOMAXPROCS(*maxProcs)
	if *size < 1 {
		fmt.Println("Size must be 1 or more")
		os.Exit(1)
	}
	c := newController(*size)
	go func() {
		if err := http.ListenAndServe(":"+*httpPort, c); err != nil {
			fmt.Printf("HTTP port listen failed %s\n", err.Error())
			os.Exit(1)
		}
	}()
	go func() {
		if err := http.ListenAndServe(":"+*profilePort, nil); err != nil {
			fmt.Printf("profile port listen failed %s\n", err.Error())
			os.Exit(1)
		}
	}()
	select {}
}

