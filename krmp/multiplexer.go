package krmp

import "io"
import "fmt"
import "log"
import "net/http"

type Multiplexer struct {
	*log.Logger
	routes     []Route
	middleware []Middleware
}

func (mux *Multiplexer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	params := make([]string, 0)
	runtime := RequestRuntime{mux.Logger, request, params}
	url := request.URL

	mux.Printf("%s %s", request.Method, request.URL.Path)

	handler := func(runtime *RequestRuntime) (Result, error) {
		return Result{}, fmt.Errorf("not found")
	}

	path := url.EscapedPath()
	for _, element := range mux.routes {
		if match := element.Path.MatchString(path) && element.Method == request.Method; match != true {
			continue
		}

		if subs := element.Path.NumSubexp(); subs == 0 {
			mux.Printf("no sub expressions found for \"%v\", moving on", element.Path)
			handler = element.Handler
			break
		}

		groups := element.Path.FindAllStringSubmatch(path, -1)

		if groups == nil {
			mux.Printf("found NO matches: %v", groups)
			handler = element.Handler
			break
		}

		for _, group := range groups[0][1:] {
			runtime.pathParams = append(runtime.pathParams, group)
		}

		mux.Printf("path params: %v", runtime.pathParams)
		handler = element.Handler
		break
	}

	for _, ware := range mux.middleware {
		handler = ware(handler)
	}

	result, err := handler(&runtime)

	if err != nil {
		mux.Printf("failed: %s", err.Error())
		writer.WriteHeader(422)
		writer.Write([]byte(err.Error()))
		return
	}

	header := writer.Header()
	header.Set("Content-Type", result.ContentType)
	header.Set("Content-Length", fmt.Sprintf("%d", result.Len()))

	if result.Attachment != "" {
		mux.Printf("adding disposition header: %s", result.Attachment)
		header.Set("Content-Disposition", fmt.Sprintf("filename=\"%s\";", result.Attachment))
	}

	writer.WriteHeader(200)
	io.Copy(writer, result)
}

func (mux *Multiplexer) Use(routes []Route, middleware []Middleware) {
	mux.middleware = middleware
	mux.routes = routes
}
