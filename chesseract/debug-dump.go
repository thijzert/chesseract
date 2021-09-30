package chesseract

import (
	"fmt"
	"io"
)

func (match Match) DebugDump(w io.Writer, highlight []Position) {
	for i, m := range match.Moves {
		if i%2 == 0 {
			fmt.Fprintf(w, " %3d: %s\n", 1+i/2, m)
		} else {
			fmt.Fprintf(w, "      %s\n", m)
		}
	}

	if _, ok := match.RuleSet.(Boring2D); ok {
		match.dumpBoring2DBoard(w, highlight)
	} else if _, ok := match.RuleSet.(Chesseract); ok {
		match.dumpHyperboard(w, highlight)
	} else {
		match.dumpUnknownBoard(w, highlight)
	}
}

func (match Match) dumpCell(w io.Writer, p Position, highlight []Position) {
	hl := false
	for _, q := range highlight {
		hl = hl || p.Equals(q)
	}

	if hl {
		if p.CellColour() == BLACK {
			w.Write([]byte("\x1b[48;5;33m\x1b[38;5;0m"))
		} else {
			w.Write([]byte("\x1b[48;5;159m\x1b[38;5;0m"))
		}
	} else {
		if p.CellColour() == BLACK {
			w.Write([]byte("\x1b[48;5;178m\x1b[38;5;0m"))
		} else {
			w.Write([]byte("\x1b[48;5;229m\x1b[38;5;0m"))
		}
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

func (match Match) dumpBoring2DBoard(w io.Writer, highlight []Position) {
	var x, y int

	fmt.Fprintf(w, "   ")
	for x = 0; x < 8; x++ {
		fmt.Fprintf(w, " %c ", 'a'+rune(x))
	}
	fmt.Fprintf(w, "\n  +------------------------+\n")
	for y = 7; y >= 0; y-- {
		fmt.Fprintf(w, "%d |", y+1)
		for x = 0; x < 8; x++ {
			match.dumpCell(w, position2D{x, y}, highlight)
		}
		fmt.Fprintf(w, "| %d\n", y+1)
	}
	fmt.Fprintf(w, "  +------------------------+\n   ")
	for x = 0; x < 8; x++ {
		fmt.Fprintf(w, " %c ", 'a'+rune(x))
	}
	fmt.Fprintf(w, "\n")
}

func (match Match) dumpHyperboard(out io.Writer, highlight []Position) {
	match.dumpUnknownBoard(out, highlight)
	var x, y, z, w int

	fmt.Fprintf(out, "   ")
	for w = 0; w < 6; w++ {
		fmt.Fprintf(out, "   ---------%d-------- ", w+1)
	}
	fmt.Fprintf(out, "\n")

	fmt.Fprintf(out, "    ")
	for w = 0; w < 6; w++ {
		fmt.Fprintf(out, "  ")
		for x = 0; x < 6; x++ {
			fmt.Fprintf(out, " %c ", 'a'+rune(x))
		}
		fmt.Fprintf(out, "  ")
	}
	fmt.Fprintf(out, "\n")
	for z = 5; z >= 0; z-- {
		fmt.Fprintf(out, "   ")
		for w = 0; w < 6; w++ {
			fmt.Fprintf(out, "  +------------------+")
		}
		fmt.Fprintf(out, "\n")
		for y = 5; y >= 0; y-- {
			if y == 3 {
				fmt.Fprintf(out, "%c| ", 'm'+rune(z))
			} else {
				fmt.Fprintf(out, " | ")
			}
			fmt.Fprintf(out, "%d", y+1)
			for w = 0; w < 6; w++ {
				fmt.Fprintf(out, " |")
				for x = 0; x < 6; x++ {
					match.dumpCell(out, position4D{x, y, z, w}, highlight)
				}
				fmt.Fprintf(out, "| ")
			}
			fmt.Fprintf(out, "%d", y+1)
			if y == 3 {
				fmt.Fprintf(out, " |%c\n", 'm'+rune(z))
			} else {
				fmt.Fprintf(out, " |\n")
			}
		}
		fmt.Fprintf(out, "   ")
		for w = 0; w < 6; w++ {
			fmt.Fprintf(out, "  +------------------+")
		}
		fmt.Fprintf(out, "\n")
	}

	fmt.Fprintf(out, "    ")
	for w = 0; w < 6; w++ {
		fmt.Fprintf(out, "  ")
		for x = 0; x < 6; x++ {
			fmt.Fprintf(out, " %c ", 'a'+rune(x))
		}
		fmt.Fprintf(out, "  ")
	}
	fmt.Fprintf(out, "\n")

	fmt.Fprintf(out, "   ")
	for w = 0; w < 6; w++ {
		fmt.Fprintf(out, "   ---------%d-------- ", w+1)
	}
	fmt.Fprintf(out, "\n")
}

func (match Match) dumpUnknownBoard(w io.Writer, highlight []Position) {
	for _, p := range match.RuleSet.AllPositions() {
		if pc, ok := match.Board.At(p); ok {
			fmt.Fprintf(w, "Position %s has %s %s\n", p, pc.Colour, pc.PieceType)
		}
	}
}
