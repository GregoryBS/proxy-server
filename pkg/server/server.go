package server

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"proxy/pkg/database"
	"strings"
	"syscall"
)

func Run(addr string) {
	db := database.GetTarantool()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		database.CloseTarantool(db)
		os.Exit(0)
	}()

	h := &Handler{db}
	server := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.handleHTTP(w, r)
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Println("Proxy started at", server.Addr)
	log.Fatal(server.ListenAndServe())
}

type Handler struct {
	db *database.DB
}

func (h *Handler) handleHTTP(w http.ResponseWriter, r *http.Request) {
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

	if err = h.db.Insert(resp, req); err != nil {
		log.Println(err)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
