package krmp

import "os"
import "io"
import "fmt"
import "log"
import "net/http"

type Multiplexer struct {
	routes     []Route
	middleware []Middleware
}

func (mux *Multiplexer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logger := log.New(os.Stdout, "krmp", log.LstdFlags)
	results := make(chan Result)
	errors := make(chan error)

	runtime := RequestRuntime{logger, request, results, errors, []string{}}
	url := request.URL

	handler := func(runtime *RequestRuntime) {
		runtime.Error(fmt.Errorf("not found"))
	}

	path := url.EscapedPath()
	for _, element := range mux.routes {
		if match := element.Path.MatchString(path) && element.Method == request.Method; match != true {
			continue
		}

		if subs := element.Path.NumSubexp(); subs == 0 {
			logger.Printf("no sub expressions found for \"%v\", moving on", element.Path)
			handler = element.Handler
			break
		}

		groups := element.Path.FindAllStringSubmatch(path, -1)

		if groups == nil {
			logger.Printf("found NO matches: %v", groups)
			handler = element.Handler
			break
		}

		for _, group := range groups[0][1:] {
			runtime.pathParams = append(runtime.pathParams, group)
		}

		logger.Printf("path params: %v", runtime.pathParams)
		handler = element.Handler
		break
	}

	for _, ware := range mux.middleware {
		handler = ware(handler)
	}

	go handler(&runtime)

	finish := func() {
		close(errors)
		close(results)
	}

	defer finish()

	select {
	case err := <-errors:
		logger.Printf("failed: %s", err.Error())
		writer.WriteHeader(422)
		writer.Write([]byte(err.Error()))
	case result := <-results:
		header := writer.Header()
		header.Set("Content-Type", result.ContentType)
		writer.WriteHeader(200)
		io.Copy(writer, result)
	}
}

func (mux *Multiplexer) Use(routes []Route, middleware []Middleware) {
	mux.middleware = middleware
	mux.routes = routes
}
