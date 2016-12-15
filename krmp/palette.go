package krmp

import "fmt"
import "github.com/lucasb-eyer/go-colorful"

type Palette struct {
	base  colorful.Color
	steps uint
	theme string
}

func (p *Palette) String() string {
	h, c, l := p.base.Hcl()
	return fmt.Sprintf("H(%f) C(%f) L(%f)", h, c, l)
}

func (p *Palette) Variations() []Variation {
	result := make([]Variation, 0)
	steps := p.steps
	head := Variation{p.base}

	amt := (360.00 / float64(steps))

	for steps > 0 {
		result = append(result, head)

		hue, saturation, brightness := head.Hsv()
		next := hue + amt

		if next >= 360 {
			next = 0.0
		}

		// move the head down some hue
		head = Variation{colorful.Hsv(next, saturation, brightness)}
		steps--
	}

	return result
}
