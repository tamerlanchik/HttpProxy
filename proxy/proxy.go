package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var (
	tunnelTimeout time.Duration
)

func Run(port string, pemPath, keyPath string, timeout int) error {
	fmt.Println("Start proxy")
	server := http.Server{
		Addr:              ":" + port,
		Handler:           http.HandlerFunc(handle),

		TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	tunnelTimeout = time.Millisecond * time.Duration(timeout)
	return server.ListenAndServeTLS(pemPath, keyPath)
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
