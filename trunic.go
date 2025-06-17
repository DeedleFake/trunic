package trunic

import (
	"bytes"
	_ "embed"
	"image"
	"io"
	"unique"

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

func makeFontMap() map[unique.Handle[string]]image.Image {
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

	r := make(map[unique.Handle[string]]image.Image, len(letters))

	for i, letter := range letters {
		x, y := i%6, i/6
		rect := image.Rect(
			letterOffsetX+x*(letterGapX+letterWidth),
			letterOffsetY+y*(letterGapY+letterHeight),
			letterOffsetX+x*(letterGapX+letterWidth)+letterWidth,
			letterOffsetY+y*(letterGapY+letterHeight)+letterHeight,
		)

		r[unique.Make(letter)] = &subimage{
			Image: fontImage,
			rect:  rect,
		}
	}

	r[unique.Make(" ")] = &subimage{
		Image: image.Transparent,
		rect:  image.Rect(0, 0, letterWidth, letterHeight),
	}

	return r
}

func validLetter(r unique.Handle[string]) bool {
	_, ok := fontMap[r]
	return ok
}

type Rune struct {
	r unique.Handle[string]
}

func MakeRune(letter string) (Rune, bool) {
	r := unique.Make(letter)
	return Rune{r: r}, validLetter(r)
}

func (r Rune) Image() image.Image {
	return fontMap[r.r]
}
