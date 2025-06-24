// Package trunic provides image-related functionality for the Trunic writing system.
package trunic

import (
	"cmp"
	"fmt"
	"maps"
	"slices"

	"deedles.dev/xiter"
	"github.com/tdewolff/canvas"
)

var (
	consonants = map[string]*canvas.Path{
		"b":  maskToPath(0b100000010100010),
		"tʃ": maskToPath(0b100000000010110),
		"d":  maskToPath(0b100000010101010),
		"f":  maskToPath(0b100000001011010),
		"ɡ":  maskToPath(0b100000001110010),
		"h":  maskToPath(0b100000010110010),
		"dʒ": maskToPath(0b100000010001010),
		"k":  maskToPath(0b100000011100010),
		"l":  maskToPath(0b100000010010010),
		"ɫ":  maskToPath(0b100000010010010),
		"m":  maskToPath(0b100000000101000),
		"n":  maskToPath(0b100000000101100),
		"ŋ":  maskToPath(0b100000011111110),
		"p":  maskToPath(0b100000001010010),
		"ɹ":  maskToPath(0b100000011010010),
		"s":  maskToPath(0b100000011011010),
		"ʃ":  maskToPath(0b100000001111110),
		"t":  maskToPath(0b100000001010110),
		"θ":  maskToPath(0b100000011010110),
		"ð":  maskToPath(0b100000010111010),
		"v":  maskToPath(0b100000010100110),
		"w":  maskToPath(0b100000001000100),
		"j":  maskToPath(0b100000010010110),
		"z":  maskToPath(0b100000010110110),
		"ʒ":  maskToPath(0b100000011101110),
	}

	vowels = map[string]*canvas.Path{
		"æ":  maskToPath(0b110011100000000),
		"ɑɹ": maskToPath(0b111100100000000),
		"ɑ":  maskToPath(0b100011100000000),
		"ɔ":  maskToPath(0b100011100000000),
		"eɪ": maskToPath(0b100000100000000),
		"ɛ":  maskToPath(0b101111000000000),
		"i":  maskToPath(0b101111100000000),
		"ɪɹ": maskToPath(0b101011100000000),
		"ə":  maskToPath(0b110000100000000),
		"ɛɹ": maskToPath(0b101011000000000),
		"ɪ":  maskToPath(0b101100000000000),
		"aɪ": maskToPath(0b110000000000000),
		"ɝ":  maskToPath(0b111111000000000),
		"oʊ": maskToPath(0b111111100000000),
		"ɔɪ": maskToPath(0b100100000000000),
		"u":  maskToPath(0b110111100000000),
		"ʊ":  maskToPath(0b100111000000000),
		"aʊ": maskToPath(0b101000000000000),
		"ɔɹ": maskToPath(0b111011100000000),
		"ʊɹ": maskToPath(0b111011100000000),
	}

	symbols = map[string]*canvas.Path{
		"*": canvas.Ellipse(.5*letterWidthRatio, .5).Translate(1, 6.5),
		".": canvas.Ellipse(.2*letterWidthRatio, .2).Translate(1, 3),
		"?": canvas.Ellipse(.1*letterWidthRatio, .1).Translate(1, 5).
			Join(canvas.EllipticalArc(.8*letterWidthRatio, .8, 0, -90, 90).Translate(1, 1.5)).
			Join(canvas.Line(0, .5).Translate(1, 1.5+1.6)),
		"!": canvas.Ellipse(.1*letterWidthRatio, .1).Translate(1, 5).
			Join(canvas.Line(0, 2).Translate(1, 1.5)),
	}

	prefixes = loadPrefixes()

	lines = []func(*canvas.Path){
		lineFunc(0, 3, 2, 3),
		lineFunc(1, 0, 2, 1),
		lineFunc(2, 5, 1, 6),
		lineFunc(1, 6, 0, 5),
		lineFunc(0, 5, 0, 4),
		lineFunc(0, 3, 0, 1),
		lineFunc(0, 1, 1, 0),
		lineFunc(1, 0, 1, 2),
		lineFunc(2, 1, 1, 2),
		lineFunc(1, 4, 2, 5),
		lineFunc(1, 4, 1, 6),
		lineFunc(1, 4, 0, 5),
		lineFunc(1, 2, 0, 1),
		lineFunc(1, 2, 1, 3),
	}
)

func maskToPath(mask int) *canvas.Path {
	var p canvas.Path
	for i, f := range lines {
		line := 1 << (len(lines) - i)
		if line&mask != 0 {
			f(&p)
		}
	}
	return &p
}

func lineFunc(x1, y1, x2, y2 float64) func(*canvas.Path) {
	return func(p *canvas.Path) {
		p.MoveTo(x1, y1)
		p.LineTo(x2, y2)
	}
}

func pathFor(ph string) *canvas.Path {
	if p := consonants[ph]; p != nil {
		return p
	}
	if p := vowels[ph]; p != nil {
		return p
	}
	if p := symbols[ph]; p != nil {
		return p
	}
	panic(fmt.Errorf("path not found for %q", ph))
}

func loadPrefixes() []string {
	return slices.SortedFunc(
		xiter.Concat(
			maps.Keys(consonants),
			maps.Keys(vowels),
			maps.Keys(symbols),
		),
		func(s1, s2 string) int { return cmp.Compare(len(s2), len(s1)) },
	)
}
