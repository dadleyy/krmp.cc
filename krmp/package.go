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
	noconflict string
}

func (p Package) rule(number int, hex, layer, modifier string) template.CSS {
	selector := fmt.Sprintf(".%s-%d", layer, number)

	if p.noconflict != "" {
		selector += fmt.Sprintf(".%s", p.noconflict)
	}

	if modifier != "" {
		selector += fmt.Sprintf(".%s", modifier)
	}

	property := "background-color"

	if layer == "fg" {
		property = "color"
	}

	rule := fmt.Sprintf("%s: %s", property, hex)

	return template.CSS(fmt.Sprintf("%s { %s; }\n", selector, rule))
}

func (p Package) Markup() template.HTML {
	result := template.HTML("<div class=\"clearfix\">")

	for i, color := range p.variations {
		result += template.HTML("<div class=\"clearfix\">")
		result += template.HTML(fmt.Sprintf("<div class=\"swatch bg-%d %s\"></div>", i, p.noconflict))

		_, _, b := color.Hsv()
		mods := map[string]int{"lighten": 1, "darken": 1}

		for _, shade := range color.Shades(p.minShade, p.maxShade, p.shadeInc) {
			_, _, sb := shade.Hsv()
			modifier := "lighten"

			if sb < b {
				modifier = "darken"
			}

			x, _ := mods[modifier]
			result += template.HTML(fmt.Sprintf("<div class=\"swatch bg-%d %s-%d %s\"></div>", i, modifier, x, p.noconflict))
			mods[modifier] += 1
		}

		result += template.HTML("</div>")
	}

	return result + template.HTML("</div>")
}

func (p Package) Stylesheet() (template.CSS, error) {
	result := template.CSS("")

	for i, color := range p.variations {
		result += p.rule(i, color.Hex(), "bg", "")
		result += p.rule(i, color.Hex(), "fg", "")

		_, _, b := color.Hsv()
		mods := map[string]int{"lighten": 1, "darken": 1}

		for _, shade := range color.Shades(p.minShade, p.maxShade, p.shadeInc) {
			_, _, sb := shade.Hsv()
			style := "lighten"

			if sb < b {
				style = "darken"
			}

			x, _ := mods[style]
			modifier := fmt.Sprintf("%s-%d", style, x)

			result += p.rule(i, shade.Hex(), "bg", modifier)
			result += p.rule(i, shade.Hex(), "fg", modifier)

			mods[style] += 1
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
