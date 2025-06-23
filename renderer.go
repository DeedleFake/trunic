package trunic

import (
	"image"
	"image/draw"
	"math"
	"slices"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
)

type Renderer struct {
	Kerning    float64
	Resolution float64

	ph [][]string
}

func (r *Renderer) resolution() canvas.Resolution {
	if r.Resolution == 0 {
		return 1
	}
	return canvas.Resolution(r.Resolution)
}

func (r *Renderer) AppendRune(ph ...string) {
	r.ph = append(r.ph, slices.Clone(ph))
}

func (r *Renderer) DrawTo(dst draw.Image) {
	resolution := r.resolution()
	renderer := rasterizer.FromImage(dst, resolution, nil)
	c := canvas.NewContext(renderer)

	for i, ph := range r.ph {
		x := float64(i) * (letterWidth + r.Kerning)

		for _, s := range ph {
			img := fontMap[s]
			c.DrawImage(x, 0, img, resolution)
		}
	}
}

func (r *Renderer) Bounds() image.Rectangle {
	return image.Rectangle{Max: r.Size()}
}

func (r *Renderer) Size() image.Point {
	return image.Pt(
		len(r.ph)*int(math.Ceil(letterWidth+r.Kerning)),
		letterHeight,
	)
}
