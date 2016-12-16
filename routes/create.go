package routes

import "bytes"
import "github.com/dadleyy/krmp.cc/krmp"

func Create(runtime *krmp.RequestRuntime) (krmp.Result, error) {
	hex := runtime.URL.Query().Get("base")

	if alt, err := runtime.PathParameter(0); err == nil {
		hex = alt
	}

	pkg, err := runtime.Package(hex)

	if err != nil {
		return krmp.Result{}, err
	}

	styles, err := pkg.Stylesheet()

	if err != nil {
		return krmp.Result{}, err
	}

	contents := bytes.NewBufferString(string(styles))
	return krmp.Result{Buffer: contents, ContentType: "text/css"}, nil
}
