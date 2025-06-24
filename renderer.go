package trunic

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"slices"
	"strings"

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

// Append is a high-level function that appends text to the Renderer's
// internal buffer. The text is expected to be a sequence of words
// written in IPA characters. Leading and trailing whitespace is
// trimmed, all unrecognized characters are stripped via [Normalize],
// the characters in words are parsed into pairs via [Runes], and
// these are then appended one-by-one via [AppendRune].
//
// If there is already anything in the buffer when Append is called,
// a space is inserted first.
//
// If the length of the text to be inserted after normalization is
// zero, this method is a no-op.
func (r *Renderer) Append(text string) {
	addSpace := func() {}
	if len(r.ph) != 0 {
		addSpace = func() { r.AppendRune() }
	}

	text = strings.TrimSpace(text)
	for word := range strings.FieldsSeq(text) {
		addSpace()
		addSpace = func() { r.AppendRune() }

		for ph := range Runes(word) {
			r.AppendRune(ph...)
		}
	}
}

// AppendRune is a low-level method that appends pieces of runes to
// the Renderer's internal buffer. Each call appends a single rune
// which is drawn by overlapping all of the symbols corresponding to
// the strings passed. Generally speaking, this is a single vowel and
// a single consonant, possibly including a reversing circle.
func (r *Renderer) AppendRune(ph ...string) {
	r.ph = append(r.ph, slices.Clone(ph))
}

// DrawTo draws the Renderer's current state to dst with the top-left
// corner at (x, y). The coordinates are relative to (0, 0), not to
// the top-left of dst's bounding box.
func (r *Renderer) DrawTo(dst draw.Image, x, y float64) {
	color := r.color()

	renderer := rasterizer.FromImage(dst, 1, nil)

	c := canvas.NewContext(renderer)
	c.SetFill(canvas.Paint{})
	c.SetStrokeColor(color)
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
		Scale(letterWidth/2, letterHeight/6.5)

	for i, ph := range r.ph {
		if len(ph) == 0 {
			continue
		}

		lx := float64(i) * letterWidth

		var p canvas.Path
		for _, ph := range ph {
			pathFor(ph).CopyTo(&p).Transform(m)
			c.DrawPath(lx, 0, &p)
		}
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

func fillPath(c *canvas.Context, x, y float64, p ...*canvas.Path) {
	stroke := c.Style.Stroke
	c.SetStroke(canvas.Paint{})
	defer c.SetStroke(stroke)

	c.SetFill(stroke)
	defer c.SetFill(canvas.Paint{})

	c.DrawPath(x, y, p...)
}
