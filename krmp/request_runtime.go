package krmp

import "log"
import "fmt"
import "bytes"
import "strings"
import "strconv"
import "net/http"
import "github.com/lucasb-eyer/go-colorful"

const PaletteDefaultBase = "#6aa7d9"
const PaletteLightenInc = 0.05
const PaletteDarkenInc = -0.05

type RequestRuntime struct {
	*log.Logger
	*http.Request
	results chan Result
	errors  chan error
}

func (runtime *RequestRuntime) Finish(r Result) {
	runtime.results <- r
}

func (runtime *RequestRuntime) Error(err error) {
	runtime.errors <- err
}

func (runtime *RequestRuntime) Palette() (Palette, error) {
	base, _ := colorful.Hex(PaletteDefaultBase)
	result := Palette{base, 36, ""}
	query := runtime.URL.Query()

	if query.Get("steps") != "" {
		steps, err := strconv.Atoi(query.Get("steps"))

		if err != nil || steps > 36 {
			return Palette{}, fmt.Errorf("step count must be valid number between 1 - 360")
		}

		result.steps = uint(steps)
	}

	if query.Get("base") != "" {
		hex := fmt.Sprintf("#%s", strings.TrimPrefix(query.Get("base"), "#"))
		base, err := colorful.Hex(hex)

		if err != nil {
			return Palette{}, fmt.Errorf("valid base color must be hexidecimal value (without leading #)")
		}

		result.base = base
	}

	return result, nil
}

func (runtime *RequestRuntime) Preview() (*bytes.Buffer, error) {
	p, err := runtime.Palette()

	if err != nil {
		return nil, err
	}

	content := ""

	for i, color := range p.Variations() {
		content += "<div class=\"clearfix\">"
		content += fmt.Sprintf("<div class=\"swatch bg-%d\"></div>", i+1)

		for x, _ := range p.Shades(color, PaletteDarkenInc) {
			content += fmt.Sprintf("<div class=\"swatch bg-%d darken-%d\"></div>", i+1, x+1)
		}

		for x, _ := range p.Shades(color, PaletteLightenInc) {
			content += fmt.Sprintf("<div class=\"swatch bg-%d lighten-%d\"></div>", i+1, x+1)
		}

		content += "</div>"
	}

	return bytes.NewBufferString(content), nil
}

func (runtime *RequestRuntime) Stylesheet() (*bytes.Buffer, error) {
	p, err := runtime.Palette()

	if err != nil {
		return nil, err
	}

	content := ""

	for i, color := range p.Variations() {
		content += fmt.Sprintf(".bg-%d { background-color: %s; }\n", i+1, color.Hex())
		content += fmt.Sprintf(".fg-%d { color: %s; }\n", i+1, color.Hex())

		for x, shade := range p.Shades(color, PaletteDarkenInc) {
			content += fmt.Sprintf(".bg-%d.darken-%d { background-color: %s; }\n", i+1, x+1, shade.Hex())
			content += fmt.Sprintf(".fg-%d.darken-%d { color: %s; }\n", i+1, x+1, shade.Hex())
		}

		for x, shade := range p.Shades(color, PaletteLightenInc) {
			content += fmt.Sprintf(".bg-%d.lighten-%d { background-color: %s; }\n", i+1, x+1, shade.Hex())
			content += fmt.Sprintf(".fg-%d.lighten-%d { color: %s; }\n", i+1, x+1, shade.Hex())
		}
	}

	return bytes.NewBufferString(content), nil
}
