package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"strings"
)

// func handleTunneling(w http.ResponseWriter, r *http.Request) {
// 	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusServiceUnavailable)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	hijacker, ok := w.(http.Hijacker)
// 	if !ok {
// 		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
// 		return
// 	}
// 	client_conn, _, err := hijacker.Hijack()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusServiceUnavailable)
// 	}
// 	go transfer(dest_conn, client_conn)
// 	go transfer(client_conn, dest_conn)
// }

// func transfer(destination io.WriteCloser, source io.ReadCloser) {
// 	defer destination.Close()
// 	defer source.Close()
// 	io.Copy(destination, source)
// }

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest(r.Method, r.RequestURI, r.Body)
	req.URL.Path = r.RequestURI[strings.Index(r.RequestURI, r.Host)+len(r.Host):]
	req.Header = r.Header.Clone()
	req.Header.Del("Proxy-Connection")
	log.Println(req)

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
	copyHeader(w.Header(), resp.Header)

}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func main() {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleHTTP(w, r)
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Println("Proxy started at", server.Addr)
	log.Fatal(server.ListenAndServe())
}
