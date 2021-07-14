package main

import (
	"fmt"
	"os"
	"time"

	"github.com/thijzert/chesseract/chesseract"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Printf("Package version: %s.  Hello, world!\n", chesseract.PackageVersion)

	rs := chesseract.Boring2D{}
	match := chesseract.Match{
		RuleSet:   rs,
		Board:     rs.DefaultBoard(),
		StartTime: time.Now(),
	}

	for {
		match.DebugDump(os.Stdout, nil)

		var move chesseract.Move
		var newBoard chesseract.Board

		for {
			fmt.Printf("Enter move for %6s: ", match.Board.Turn)

			var sFrom, sTo string
			n, _ := fmt.Scanf("%s %s\n", &sFrom, &sTo)
			if n == 0 {
				continue
			}
			if n == 1 {
				if sFrom == "forfeit" || sFrom == "quit" {
					fmt.Printf("%s forfeits", match.Board.Turn)
					return nil
				}
			}

			from, err := rs.ParsePosition(sFrom)
			if err != nil {
				fmt.Printf("error parsing '%s': %v\n", sFrom, err)
				continue
			}
			piece, _ := match.Board.At(from)
			to, err := rs.ParsePosition(sTo)
			if err != nil {
				fmt.Printf("error parsing '%s': %v\n", sTo, err)
				continue
			}

			moveTime := time.Since(match.StartTime)
			for _, m := range match.Moves {
				moveTime -= m.Time
			}

			move = chesseract.Move{
				PieceType: piece.PieceType,
				From:      from,
				To:        to,
				Time:      moveTime,
			}
			newBoard, err = rs.ApplyMove(match.Board, move)
			if err != nil {
				fmt.Printf("applying move '%s'-'%s': %v\n", sFrom, sTo, err)
				continue
			}

			break
		}

		match.Moves = append(match.Moves, move)
		match.Board = newBoard
	}
}
