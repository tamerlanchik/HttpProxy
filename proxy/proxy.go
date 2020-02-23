package proxy

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	DB "proxy/db"
	"strings"
	"time"
)

var (
	tunnelTimeout time.Duration
	db *sql.DB
)
const(
	tmpl = "|%7s|%8s|%30s|%20s|%50s|\n"
)
func Run(port string, pemPath, keyPath string, timeout int, proto string) error {
	fmt.Println("Start proxy")
	fmt.Printf("Port: %s, Proto: %s\n", port, proto)
	fmt.Printf(tmpl, "Meth", "Sche", "Destination", "Created", "Header")
	server := http.Server{
		Addr:              ":" + port,
		Handler:           http.HandlerFunc(handle),

		TLSNextProto:      make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	tunnelTimeout = time.Millisecond * time.Duration(timeout)

	var err error
	db, err = DB.Connect()
	if err != nil {
		return err
	}
	if proto == "http" {
		return server.ListenAndServe()
	} else {
		return server.ListenAndServeTLS(pemPath, keyPath)
	}
}

func handle(w http.ResponseWriter, r* http.Request) {
	//fmt.Print(r.URL.Path[1:] + ": ")
	if r.Method == http.MethodConnect {
		HandleTunnel(w, r)
	} else {
		HandleHttp(w, r)
	}
	query := `INSERT INTO Request(method, proto_schema, dest, body, header, created) 
				VALUES ($1, $2, $3, $4, $5, $6);`
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	var headers string
	for key, values := range r.Header {
		headers += key + ": " + strings.Join(values, ", ") + "\n"
	}
	_, err = db.Exec(query, r.Method, r.Proto, r.Host + r.URL.Path, string(body),
		headers, time.Now().Format(time.RFC822))
	if err != nil {
		fmt.Printf("Cannot write to database: %s\n", err.Error())
	}
	fmt.Printf(tmpl, r.Method, r.Proto, r.Host+r.URL.Path, time.Now().Format(time.RFC822), "\n"+headers)
}

func HandleHttp(w http.ResponseWriter, r* http.Request) {
	//fmt.Println("Http")
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
	//fmt.Println("Https")
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
