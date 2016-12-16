package routes

import "fmt"
import "bytes"
import "archive/tar"
import "compress/gzip"
import "github.com/dadleyy/krmp.cc/krmp"

func Download(runtime *krmp.RequestRuntime) (krmp.Result, error) {
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

	bower, err := pkg.Bowerfile()

	if err != nil {
		return krmp.Result{}, err
	}

	buffer := bytes.NewBuffer(make([]byte, 0))

	compressor := gzip.NewWriter(buffer)
	archiver := tar.NewWriter(compressor)

	defer compressor.Close()
	defer archiver.Close()

	header := &tar.Header{
		Name: "krmp/bower.json",
		Mode: 0755,
		Size: int64(len(bower)),
	}

	if err := archiver.WriteHeader(header); err != nil {
		return krmp.Result{}, err
	}

	if _, err := archiver.Write([]byte(bower)); err != nil {
		return krmp.Result{}, err
	}

	header = &tar.Header{
		Name: fmt.Sprintf("krmp/base-%s.css", hex),
		Mode: 0755,
		Size: int64(len(styles)),
	}

	if err := archiver.WriteHeader(header); err != nil {
		return krmp.Result{}, err
	}

	if _, err := archiver.Write([]byte(styles)); err != nil {
		return krmp.Result{}, err
	}

	return krmp.Result{buffer, "application/x-gzip", fmt.Sprintf("krmp-%s.tar.gz", hex)}, nil
}
