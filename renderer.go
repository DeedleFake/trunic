package trunic

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"slices"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
)

const letterWidthRatio = .8

type Renderer struct {
	Kerning    float64     // Default: 0
	TextHeight float64     // Default: 24
	Color      color.Color // Default: color.Black
	Thickness  float64     // Default: 10
	Resolution float64     // Default: 3

	ph [][]string
}

func (r *Renderer) textHeight() float64 {
	if r.TextHeight == 0 {
		return 24
	}
	return r.TextHeight
}

func (r *Renderer) resolution() canvas.Resolution {
	if r.Resolution == 0 {
		return 3
	}
	return canvas.Resolution(r.Resolution)
}

func (r *Renderer) color() color.Color {
	if r.Color == nil {
		return color.Black
	}
	return r.Color
}

func (r *Renderer) thickness() float64 {
	if r.Thickness == 0 {
		return 10
	}
	return r.Thickness
}

func (r *Renderer) AppendRune(ph ...string) {
	r.ph = append(r.ph, slices.Clone(ph))
}

func (r *Renderer) DrawTo(dst draw.Image, x, y float64) {
	resolution := r.resolution()
	offset := dst.Bounds().Canon().Min
	x -= float64(offset.X) / float64(resolution)
	y -= float64(offset.Y) / float64(resolution)

	renderer := rasterizer.FromImage(dst, resolution, nil)

	c := canvas.NewContext(renderer)
	c.SetStrokeColor(r.color())
	c.SetStrokeWidth(r.thickness())
	c.SetStrokeCapper(canvas.RoundCap)
	c.SetCoordSystem(canvas.CartesianIV)

	letterHeight := r.textHeight()
	letterWidth := letterWidthRatio * letterHeight

	for i := range r.ph {
		lx := x + float64(i)*(letterWidth+r.Kerning)

		c.MoveTo(lx, y)
		c.LineTo(lx+letterWidth, y+letterHeight)
		c.Stroke()
	}
}

func (r *Renderer) Bounds() image.Rectangle {
	return image.Rectangle{Max: r.Size()}
}

func (r *Renderer) Size() image.Point {
	height := float64(r.resolution()) * r.textHeight()
	width := letterWidthRatio * height

	return image.Pt(
		len(r.ph)*int(math.Ceil(width+r.Kerning)),
		int(math.Ceil(height)),
	)
}
