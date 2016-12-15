package krmp

import "log"
import "fmt"
import "strings"
import "strconv"
import "net/http"
import "github.com/lucasb-eyer/go-colorful"

const PaletteDefaultBase = "#6aa7d9"
const PaletteDefaultSteps = 3
const PaletteDefaultShadeInc = 0.05
const PaletteDefaultShadeMax = 0.99
const PaletteDefaultShadeMin = 0.50

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

func (runtime *RequestRuntime) Package() (Package, error) {
	base, _ := colorful.Hex(PaletteDefaultBase)
	result := Palette{base, PaletteDefaultSteps, ""}
	query := runtime.URL.Query()

	if query.Get("steps") != "" {
		steps, err := strconv.Atoi(query.Get("steps"))

		if err != nil || steps > 36 {
			return Package{}, fmt.Errorf("step count must be valid number between 1 - 360")
		}

		result.steps = uint(steps)
	}

	if query.Get("base") != "" {
		hex := fmt.Sprintf("#%s", strings.TrimPrefix(query.Get("base"), "#"))
		base, err := colorful.Hex(hex)

		if err != nil {
			return Package{}, fmt.Errorf("valid base color must be hexidecimal value (without leading #)")
		}

		result.base = base
	}

	min, max, inc := PaletteDefaultShadeMin, PaletteDefaultShadeMax, PaletteDefaultShadeInc

	if query.Get("shade_max") != "" {
		m, err := strconv.Atoi(query.Get("shade_max"))

		if err != nil || m > 100 || m < 0 {
			return Package{}, fmt.Errorf("shade_max must be between 0 - 100")
		}

		min = float64(m) * 0.01
	}

	if query.Get("shade_min") != "" {
		m, err := strconv.Atoi(query.Get("shade_min"))

		if err != nil || m > 100 || m < 0 {
			return Package{}, fmt.Errorf("shade_min must be between 0 - 100")
		}

		min = float64(m) * 0.01
	}

	if query.Get("shades") == "false" {
		min = 100
		max = 0
	}

	expanded := query.Get("expanded") == "true"
	return Package{result.Variations(), min, max, inc, expanded}, nil
}
