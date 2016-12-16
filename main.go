package main

import "os"
import "fmt"
import "log"
import "flag"
import "time"
import "regexp"
import "net/http"

import "github.com/dadleyy/krmp.cc/krmp"
import "github.com/dadleyy/krmp.cc/routes"

func main() {
	var port string

	flag.StringVar(&port, "port", "8080", "which port to run the server on")
	flag.Parse()

	logger := log.New(os.Stdout, "[krmp.cc] ", log.LstdFlags)
	logger.Printf("starting server on port %s\n", port)

	mux := krmp.Multiplexer{Logger: logger}

	routes := []krmp.Route{
		krmp.Route{"GET", regexp.MustCompile("^/preview$"), routes.Preview},
		krmp.Route{"GET", regexp.MustCompile("^/([a-f0-9]{6}|[a-f0-9]{3})$"), routes.Create},
		krmp.Route{"GET", regexp.MustCompile("^/([a-f0-9]{6}|[a-f0-9]{3})/preview$"), routes.Preview},
		krmp.Route{"GET", regexp.MustCompile("^/([a-f0-9]{6}|[a-f0-9]{3})/download$"), routes.Download},
		krmp.Route{"GET", regexp.MustCompile(".*"), routes.Create},
	}

	middleware := []krmp.Middleware{}

	mux.Use(routes, middleware)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        &mux,
	}

	server.ListenAndServe()
}
