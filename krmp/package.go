package krmp

import "fmt"
import "bytes"
import "html/template"
import "encoding/json"
import "github.com/tdewolff/minify"
import "github.com/tdewolff/minify/css"

const PackageDefaultName = "cc"

type Package struct {
	variations []Variation
	minShade   float64
	maxShade   float64
	shadeInc   float64
	expanded   bool
	noconflict string
	rules      map[string]string
}

func (p Package) Markup() template.HTML {
	result := template.HTML("<div class=\"clearfix\">")

	bg, ok := p.rules["background-color"]

	if ok != true {
		return template.HTML("")
	}

	for i, color := range p.variations {
		result += template.HTML("<div class=\"clearfix\">")
		result += template.HTML(fmt.Sprintf("<div class=\"swatch %s-%d %s\"></div>", bg, i, p.noconflict))

		_, _, b := color.Hsv()
		mods := map[string]int{"lighten": 1, "darken": 1}

		for _, shade := range color.Shades(p.minShade, p.maxShade, p.shadeInc) {
			_, _, sb := shade.Hsv()
			modifier := "lighten"

			if sb < b {
				modifier = "darken"
			}

			x, _ := mods[modifier]
			result += template.HTML(fmt.Sprintf("<div class=\"swatch %s-%d %s-%d %s\"></div>", bg, i, modifier, x, p.noconflict))
			mods[modifier] += 1
		}

		result += template.HTML("</div>")
	}

	return result + template.HTML("</div>")
}

func (p Package) Bowerfile() (template.JS, error) {
	name := fmt.Sprintf("krmp-%s", PackageDefaultName)

	if p.noconflict != "" {
		name = fmt.Sprintf("krmp-%s", p.noconflict)
	}

	definition := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}{name, "on the fly css pallete stylsheet"}

	data, err := json.Marshal(definition)

	if err != nil {
		return template.JS(""), err
	}

	return template.JS(string(data)), nil
}

func (p Package) Stylesheet() (template.CSS, error) {
	result := template.CSS("")

	for i, color := range p.variations {
		for name, alias := range p.rules {
			class := fmt.Sprintf("%s-%d", alias, i)
			result += template.CSS(fmt.Sprintf(".%s { %s: %s; }\n", class, name, color.Hex()))
		}

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

			for name, alias := range p.rules {
				class := fmt.Sprintf("%s-%d", alias, i)
				result += template.CSS(fmt.Sprintf(".%s.%s { %s: %s; }\n", class, modifier, name, shade.Hex()))
			}

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
