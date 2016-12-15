package main

import "fmt"
import "time"
import "bytes"
import "regexp"
import "net/http"
import "html/template"

import "github.com/dadleyy/krmp.cc/krmp"

func logger(next krmp.Terminal) krmp.Terminal {
	exec := func(runtime *krmp.RequestRuntime) {
		runtime.Printf("%s %s %s", runtime.Request.Method, runtime.URL.EscapedPath(), runtime.URL.RawQuery)
		next(runtime)
	}

	return exec
}

func create(runtime *krmp.RequestRuntime) {
	stylesheet, err := runtime.Stylesheet()

	if err != nil {
		runtime.Error(err)
		return
	}

	runtime.Finish(krmp.Result{stylesheet, "text/css"})
}

func preview(runtime *krmp.RequestRuntime) {
	stylesheet, err := runtime.Stylesheet()

	if err != nil {
		runtime.Error(err)
		return
	}

	markup, err := runtime.Preview()

	if err != nil {
		runtime.Error(err)
		return
	}

	engine, err := template.ParseFiles("preview.html")

	if err != nil {
		runtime.Error(err)
		return
	}

	buffer := bytes.NewBuffer(make([]byte, 0))

	context := struct {
		Styles   template.CSS
		Previews template.HTML
	}{template.CSS(stylesheet.String()), template.HTML(markup.String())}

	if err := engine.Execute(buffer, context); err != nil {
		runtime.Error(err)
		return
	}

	runtime.Finish(krmp.Result{buffer, "text/html"})
}

func main() {
	fmt.Printf("starting server...\n")
	port := "8080"

	mux := krmp.Multiplexer{}

	routes := []krmp.Route{
		krmp.Route{"GET", regexp.MustCompile("/preview"), preview},
		krmp.Route{"GET", regexp.MustCompile(".*"), create},
	}

	middleware := []krmp.Middleware{logger}

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
