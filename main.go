package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var (
	tunnelTimeout time.Duration
)

//	Test:
//	curl -Lv --proxy https://localhost:12345 --proxy-cacert server.pem https://google.com
func main() {
	var pemPath, keyPath string
	var tunnelTimeoutTemp int64
	flag.StringVar(&pemPath, "pem", "keys/server.pem", "Path to pem-file")
	flag.StringVar(&keyPath, "key", "keys/server.key", "Path to key-file")
	flag.Int64Var(&tunnelTimeoutTemp, "timeout", 1000, "Timeout, ms")
	flag.Parse()
	tunnelTimeout = tunnelTimeout * time.Millisecond

	fmt.Println("Hello")
	server := http.Server{
		Addr:              ":12345",
		Handler:           http.HandlerFunc(handle),

		TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	//http.HandleFunc("/", handle);
	//fmt.Println(http.ListenAndServe(":8080", nil))
	fmt.Println(server.ListenAndServeTLS(pemPath, keyPath))

}

func handle(w http.ResponseWriter, r* http.Request) {
	//fmt.Print(r.URL.Path[1:] + ": ")
	if r.Method == http.MethodConnect {
		HandleTunnel(w, r)
	} else {
		HandleHttp(w, r)
	}
}

func HandleHttp(w http.ResponseWriter, r* http.Request) {
	fmt.Println("Http")
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer r.Body.Close()

	copyHeader := func(dest, src http.Header) {
		for key, valueSlice := range src {
			for _, value := range valueSlice {
				dest.Add(key, value)
			}
		}
	}
	copyHeader(w.Header(), r.Header)
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func HandleTunnel(w http.ResponseWriter, r* http.Request) {
	fmt.Println("Https")
	// Dest connection
	destinationConnection, err := net.DialTimeout("tcp", r.Host, tunnelTimeout)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable);
		return
	}

	// Source connection
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking is not supported", http.StatusInternalServerError)
		return
	}

	clientConnection, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	// make Tunnel
	wormHole := func(dest io.WriteCloser, src io.ReadCloser) {
		defer dest.Close()
		defer src.Close()
		io.Copy(dest, src)
	}
	go wormHole(destinationConnection, clientConnection)
	go wormHole(clientConnection, destinationConnection)
}
