// Package trunic provides image-related functionality for the Trunic writing system.
package trunic

import "github.com/tdewolff/canvas"

var base = canvas.MustParseSVGPath("M0 3 L2 3")

var runes = map[string]*canvas.Path{
	"t": canvas.MustParseSVGPath("M0 1 L1 2 L2 1 M1 2 L1 3 M1 4 L1 6"),
	"e": canvas.MustParseSVGPath("M0 1 L0 3 M0 4 L0 5 L1 6 L2 5"),
	"s": canvas.MustParseSVGPath("M1 0 L1 3 L2 1 M0 5 L1 4 L1 6"),
}
