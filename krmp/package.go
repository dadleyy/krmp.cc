package krmp

import "fmt"
import "bytes"
import "html/template"
import "github.com/tdewolff/minify"
import "github.com/tdewolff/minify/css"

type Package struct {
	variations []Variation
	minShade   float64
	maxShade   float64
	shadeInc   float64
	expanded   bool
}

func (p Package) Markup() template.HTML {
	result := template.HTML("<div class=\"clearfix\">")

	for i, color := range p.variations {
		result += template.HTML("<div class=\"clearfix\">")
		result += template.HTML(fmt.Sprintf("<div class=\"swatch bg-%d\"></div>", i))

		_, _, b := color.Hsv()
		mods := map[string]int{"lighten": 1, "darken": 1}

		for _, shade := range color.Shades(p.minShade, p.maxShade, p.shadeInc) {
			_, _, sb := shade.Hsv()
			modifier := "lighten"

			if sb < b {
				modifier = "darken"
			}

			x, _ := mods[modifier]
			result += template.HTML(fmt.Sprintf("<div class=\"swatch bg-%d %s-%d\"></div>", i, modifier, x))
			mods[modifier] += 1
		}

		result += template.HTML("</div>")
	}

	return result + template.HTML("</div>")
}

func (p Package) Stylesheet() (template.CSS, error) {
	result := template.CSS("")

	for i, color := range p.variations {
		result += template.CSS(fmt.Sprintf(".bg-%d { background-color: %s; }\n", i, color.Hex()))
		result += template.CSS(fmt.Sprintf(".fg-%d { color: %s; }\n", i, color.Hex()))

		_, _, b := color.Hsv()
		mods := map[string]int{"lighten": 1, "darken": 1}

		for _, shade := range color.Shades(p.minShade, p.maxShade, p.shadeInc) {
			_, _, sb := shade.Hsv()
			modifier := "lighten"

			if sb < b {
				modifier = "darken"
			}

			x, _ := mods[modifier]
			result += template.CSS(fmt.Sprintf(".bg-%d.%s-%d { background-color: %s; }\n", i, modifier, x, shade.Hex()))
			mods[modifier] += 1
		}
	}

	if p.expanded == true {
		return result, nil
	}

	compressor := minify.New()
	compressor.AddFunc("text/css", css.Minify)

	input := bytes.NewBufferString(string(result))
	output := bytes.NewBuffer([]byte{})

	if err := compressor.Minify("text/css", output, input); err != nil {
		return template.CSS(""), err
	}

	return template.CSS(output.String()), nil
}
