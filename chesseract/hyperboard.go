package chesseract

import "fmt"

// A position4D represents a position on the hyper-board
type position4D [4]int

func (p position4D) String() string {
	return fmt.Sprintf("%c%d%c%d", 'a'+rune(p[0]), p[1]+1, 'm'+rune(p[2]), p[3]+1)
}

func (p position4D) Equals(q Position) bool {
	if q0, ok := q.(position4D); ok {
		return p == q0
	}
	return false
}

func (p position4D) CellColour() Colour {
	if (p[0]+p[1]+p[2]+p[3])%2 == 0 {
		return BLACK
	} else {
		return WHITE
	}
}
