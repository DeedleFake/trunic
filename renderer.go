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
	Color      color.Color // Default: color.Black
	TextHeight float64     // Default: 72
	Thickness  float64     // Default: 10

	ph [][]string
}

func (r *Renderer) color() color.Color {
	if r.Color == nil {
		return color.Black
	}
	return r.Color
}

func (r *Renderer) textHeight() float64 {
	if r.TextHeight == 0 {
		return 72
	}
	return r.TextHeight
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
	offset := dst.Bounds().Canon().Min
	ox, oy := float64(-offset.X), float64(-offset.Y)

	renderer := rasterizer.FromImage(dst, 1, nil)

	c := canvas.NewContext(renderer)
	c.SetStrokeColor(r.color())
	c.SetStrokeWidth(r.thickness())
	c.SetStrokeCapper(canvas.RoundCap)
	c.SetCoordSystem(canvas.CartesianIV)

	letterHeight := r.textHeight()
	letterWidth := letterWidthRatio * letterHeight

	var py float64
	for i := range r.ph {
		lx := x + float64(i)*letterWidth
		m := canvas.Identity.Translate(ox, oy).Scale(letterWidth, letterHeight)

		var p canvas.Path
		p.MoveTo(0, py)
		p.LineTo(1, 1-py)
		c.DrawPath(lx, 0, p.Transform(m))
		c.Stroke()

		py = 1 - py
	}
}

func (r *Renderer) Bounds() image.Rectangle {
	return image.Rectangle{Max: r.Size()}
}

func (r *Renderer) Size() image.Point {
	height := r.textHeight()
	width := letterWidthRatio * height

	return image.Pt(
		len(r.ph)*int(math.Ceil(width)),
		int(math.Ceil(height)),
	)
}

func t(x, y float64, m canvas.Matrix) (float64, float64) {
	p := m.Dot(canvas.Point{X: x, Y: y})
	return p.X, p.Y
}
