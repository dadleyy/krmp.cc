package krmp

import "github.com/lucasb-eyer/go-colorful"

type Variation struct {
	colorful.Color
}

func (v Variation) Shades(min, max, inc float64) []Variation {
	result := make([]Variation, 0)

	hue, saturation, brightness := v.Hsv()

	for b := brightness + inc; b < max; b += inc {
		result = append(result, Variation{colorful.Hsv(hue, saturation, b)})
	}

	for b := brightness - inc; b > min; b -= inc {
		result = append(result, Variation{colorful.Hsv(hue, saturation, b)})
	}

	return result
}
