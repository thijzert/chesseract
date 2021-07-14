package chesseract

import (
	"fmt"
	"io"
)

func (match Match) DebugDump(w io.Writer) {
	if _, ok := match.RuleSet.(Boring2D); ok {
		match.dumpBoring2DBoard(w)
	} else {
		match.dumpUnknownBoard(w)
	}

	for i, m := range match.Moves {
		if i%2 == 0 {
			fmt.Fprintf(w, " %3d: %s\n", 1+i/2, m)
		} else {
			fmt.Fprintf(w, "      %s\n", m)
		}
	}
}

func (match Match) dumpCell(w io.Writer, p Position) {
	if p.CellColour() == BLACK {
		w.Write([]byte("\x1b[48;5;178m\x1b[38;5;0m"))
	} else {
		w.Write([]byte("\x1b[48;5;229m\x1b[38;5;0m"))
	}

	if pc, ok := match.Board.At(p); ok {
		s := pc.PieceType.String()
		var r rune
		for _, rr := range s {
			r = rr
			break
		}
		if pc.Colour == WHITE && r > 256 {
			r -= 6
		}
		fmt.Fprintf(w, " %c ", r)
	} else {
		w.Write([]byte{32, 32, 32})
	}

	w.Write([]byte("\x1b[0m"))
}

func (match Match) dumpBoring2DBoard(w io.Writer) {
	var x, y int

	fmt.Fprintf(w, "   ")
	for x = 0; x < 8; x++ {
		fmt.Fprintf(w, " %c ", 'a'+rune(x))
	}
	fmt.Fprintf(w, "\n  +------------------------+\n")
	for y = 7; y >= 0; y-- {
		fmt.Fprintf(w, "%d |", y+1)
		for x = 0; x < 8; x++ {
			match.dumpCell(w, position2D{x, y})
		}
		fmt.Fprintf(w, "| %d\n", y+1)
	}
	fmt.Fprintf(w, "  +------------------------+\n   ")
	for x = 0; x < 8; x++ {
		fmt.Fprintf(w, " %c ", 'a'+rune(x))
	}
	fmt.Fprintf(w, "\n")
}

func (match Match) dumpUnknownBoard(w io.Writer) {
	for _, p := range match.RuleSet.AllPositions() {
		if pc, ok := match.Board.At(p); ok {
			fmt.Fprintf(w, "Position %s has %s %s\n", p, pc.Colour, pc.PieceType)
		}
	}
}
