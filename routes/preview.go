package routes

import "bytes"
import "html/template"
import "github.com/dadleyy/krmp.cc/krmp"

func Preview(runtime *krmp.RequestRuntime) (krmp.Result, error) {
	hex := runtime.URL.Query().Get("base")

	if alt, err := runtime.PathParameter(0); err == nil {
		hex = alt
	}

	runtime.Printf("previewing hex \"%s\"", hex)
	pkg, err := runtime.Package(hex)

	if err != nil {
		return krmp.Result{}, err
	}

	engine, err := template.ParseFiles("preview.html")

	if err != nil {
		return krmp.Result{}, err
	}

	buffer := bytes.NewBuffer(make([]byte, 0))

	styles, err := pkg.Stylesheet()

	if err != nil {
		return krmp.Result{}, err
	}

	context := struct {
		Styles   template.CSS
		Previews template.HTML
	}{styles, pkg.Markup()}

	if err := engine.Execute(buffer, context); err != nil {
		return krmp.Result{}, err
	}

	return krmp.Result{Buffer: buffer, ContentType: "text/html"}, nil
}
