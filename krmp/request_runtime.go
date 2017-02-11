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
const PaletteMaxSteps = 36
const PaletteMaxNoConflictSize = 6

type RequestRuntime struct {
	*log.Logger
	*http.Request
	pathParams []string
}

func destructure(in string, on string) (string, string, error) {
	parts := strings.Split(in, on)

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid string split")
	}

	return parts[0], parts[1], nil
}

func (runtime *RequestRuntime) PathParameter(index int) (string, error) {
	if runtime.pathParams == nil || index >= len(runtime.pathParams) {
		return "", fmt.Errorf("parameter not found")
	}

	return runtime.pathParams[index], nil
}

func (runtime *RequestRuntime) Package(hex string) (Package, error) {
	base, _ := colorful.Hex(PaletteDefaultBase)
	result := Palette{base, PaletteDefaultSteps, ""}
	query := runtime.URL.Query()

	if query.Get("steps") != "" {
		steps, err := strconv.Atoi(query.Get("steps"))

		if err != nil || steps > PaletteMaxSteps {
			return Package{}, fmt.Errorf("step count must be valid number between 1 - %d", PaletteMaxSteps)
		}

		result.steps = uint(steps)
	}

	if hex != "" {
		hex = fmt.Sprintf("#%s", strings.TrimPrefix(hex, "#"))
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

		max = float64(m) * 0.01
	}

	if query.Get("shade_min") != "" {
		m, err := strconv.Atoi(query.Get("shade_min"))

		if err != nil || m > 100 || m < 0 {
			return Package{}, fmt.Errorf("shade_min must be between 0 - 100")
		}

		min = float64(m) * 0.01
	}

	rules := map[string]string{
		"background-color": "bg",
		"color":            "fg",
		"border-color":     "bc",
	}

	if val, ok := query["rules"]; ok && len(val) == 1 {
		list := strings.Split(val[0], ",")

		for _, mapping := range list {
			name, alt, err := destructure(mapping, ":")

			if err != nil {
				continue
			}

			if _, exists := rules[name]; exists != true {
				continue
			}

			if alt == "-" || alt == "" {
				delete(rules, name)
				continue
			}

			rules[name] = alt
		}
	}

	fmt.Printf("current rules: %v\n", rules)

	if query.Get("shades") == "false" {
		min = 100
		max = 0
	}

	if alt := query.Get("noconflict"); len(alt) > PaletteMaxNoConflictSize {
		return Package{}, fmt.Errorf("noconflict classes must be less than 7 characters long")
	}

	return Package{
		variations: result.Variations(),
		minShade:   min,
		maxShade:   max,
		shadeInc:   inc,
		expanded:   query.Get("expanded") == "true",
		noconflict: query.Get("noconflict"),
		rules:      rules,
	}, nil
}
