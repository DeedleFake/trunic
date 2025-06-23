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

const letterWidthRatio = .6

// Renderer draws Trunic text to an image.
type Renderer struct {
	Color      color.Color // Default: color.Black
	TextHeight float64     // Default: 72
	Thickness  float64     // Default: 5

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
		return 5
	}
	return r.Thickness
}

// AppendRune is a low-level method that appends pieces of runes to
// the Renderer's internal buffer. Each call appends a single rune
// which is drawn by overlapping all of the symbols corresponding to
// the strings passed. Generally speaking, this is a single vowel and
// a single consonant.
func (r *Renderer) AppendRune(ph ...string) {
	r.ph = append(r.ph, slices.Clone(ph))
}

// DrawTo draws the Renderer's current state to dst with the top-left
// corner at (x, y). The coordinates are relative to (0, 0), not to
// the top-left of dst's bounding box.
func (r *Renderer) DrawTo(dst draw.Image, x, y float64) {
	renderer := rasterizer.FromImage(dst, 1, nil)

	c := canvas.NewContext(renderer)
	c.SetFill(canvas.Paint{})
	c.SetStrokeColor(r.color())
	c.SetStrokeWidth(r.thickness())
	c.SetStrokeCapper(canvas.RoundCap)
	c.SetCoordSystem(canvas.CartesianIV)
	c.SetStrokeJoiner(canvas.RoundJoin)

	letterHeight := r.textHeight()
	letterWidth := letterWidthRatio * letterHeight

	offset := dst.Bounds().Canon().Min
	m := canvas.Identity.
		Translate(x, y).
		Translate(float64(-offset.X), float64(-offset.Y)).
		Scale(letterWidth/2, letterHeight/6)

	for i, ph := range r.ph {
		if len(ph) == 0 {
			continue
		}

		p := &canvas.Path{}
		for _, ph := range ph {
			p = p.Join(pathFor(ph))
		}

		lx := float64(i) * letterWidth
		c.DrawPath(lx, 0, p.Transform(m))
	}
}

// Bounds returns the minimum bounding box that will contain the
// result of rendering the Renderer's current state with the top-left
// corner at (0, 0). Note that this is a minimum that does not take
// into account the thickness of the lines. Default thickness usually
// requires about 10 pixels of padding or so to fit the result without
// anything getting clipped.
func (r *Renderer) Bounds() image.Rectangle {
	return image.Rectangle{Max: r.Size()}
}

// Size returns the size of the image. This is equivalent to
//
//	r.Bounds().Size()
func (r *Renderer) Size() image.Point {
	height := r.textHeight()
	width := letterWidthRatio * height

	return image.Pt(
		len(r.ph)*int(math.Ceil(width)),
		int(math.Ceil(height)),
	)
}
