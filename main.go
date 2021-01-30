package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"text/template"
	"time"
)

func main() {
	listenAddr := "8080"
	if len(os.Args) == 2 {
		listenAddr = os.Args[1]
	}

	// create a logger, router and server
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	router := http.NewServeMux()
	router.HandleFunc("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))).ServeHTTP)
	router.HandleFunc("/", index)
	router.HandleFunc("/healthz", healthz)
	router.HandleFunc("/search", searchHandler)
	router.HandleFunc("/detail", detailHandler)

	server := newServer(
		listenAddr,
		(middlewares{
			logging(logger),
			tracing(func() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }),
		}).apply(router),
		logger,
	)

	// run our server
	if err := server.run(); err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, nil)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(w, `{"alive": true}`)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}
