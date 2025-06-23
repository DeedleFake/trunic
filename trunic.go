// Package trunic provides image-related functionality for the Trunic writing system.
package trunic

import (
	"bytes"
	_ "embed"
	"image"
	"image/draw"
	"io"

	"golang.org/x/image/webp"
)

const (
	letterOffsetX = 2
	letterOffsetY = 0
	letterWidth   = 78
	letterHeight  = 122
	letterGapX    = 32
	letterGapY    = 56
)

var (
	//go:embed trunic.webp
	fontImageBytes []byte
	fontImage      = loadImage(bytes.NewReader(fontImageBytes))

	fontMap = makeFontMap()
)

func loadImage(r io.Reader) image.Image {
	img, err := webp.Decode(r)
	if err != nil {
		panic(err)
	}
	return img
}

func makeFontMap() map[string]image.Image {
	letters := [...]string{
		"a", "ar", "ah", "ay", "e", "ee",
		"eer", "u", "er", "i", "ie", "ir",
		"o", "oy", "oo", "ou", "ow", "or",
		"b", "ch", "d", "f", "g", "h",
		"j", "k", "l", "m", "n", "ng",
		"p", "r", "s", "sh", "t", "th",
		"dh", "v", "w", "y", "z", "zh",
		"-", ",", ".", "!", "?", "_",
	}

	r := make(map[string]image.Image, len(letters))

	for i, letter := range letters {
		x, y := i%6, i/6
		rect := image.Rect(
			letterOffsetX+x*(letterGapX+letterWidth),
			letterOffsetY+y*(letterGapY+letterHeight),
			letterOffsetX+x*(letterGapX+letterWidth)+letterWidth,
			letterOffsetY+y*(letterGapY+letterHeight)+letterHeight,
		)

		r[letter] = &subimage{
			Image: fontImage,
			rect:  rect,
		}
	}

	r[" "] = &subimage{
		Image: image.Transparent,
		rect:  image.Rect(0, 0, letterWidth, letterHeight),
	}

	return r
}

type Renderer struct {
	ph [][]image.Image
}

func (r *Renderer) AppendRune(ph ...string) {
	imgs := make([]image.Image, 0, len(ph))
	for _, ph := range ph {
		imgs = append(imgs, fontMap[ph])
	}

	r.ph = append(r.ph, imgs)
}

func (r *Renderer) Draw(dst draw.Image, dr image.Rectangle, src image.Image, sp image.Point) {
	dr = dr.Canon()

	for i, ph := range r.ph {
		dp := dr.Min.Add(image.Pt(i*letterWidth, 0))
		dr := image.Rectangle{Min: dp, Max: dp.Add(image.Pt(letterWidth, letterHeight))}.Intersect(dr)

		for _, img := range ph {
			draw.DrawMask(
				dst,
				dr,
				src,
				sp,
				img,
				image.Point{},
				draw.Over,
			)
		}
	}
}

func (r *Renderer) Bounds() image.Rectangle {
	return image.Rectangle{Max: r.Size()}
}

func (r *Renderer) Size() image.Point {
	return image.Pt(
		len(r.ph)*letterWidth,
		letterHeight,
	)
}
