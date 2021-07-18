package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/thijzert/chesseract/chesseract"
	plumbing "github.com/thijzert/chesseract/internal/web-plumbing"
)

func main() {
	fmt.Printf("Chesseract version: %s\n", chesseract.PackageVersion)

	var err error = fmt.Errorf("invalid command")
	if len(os.Args) < 2 {
		err = consoleGame()
	} else if os.Args[1] == "server" {
		err = apiServer()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func consoleGame() error {
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

func apiServer() error {
	var listenPort string
	var storageBackend string

	fs := flag.NewFlagSet(os.Args[0]+" server", flag.ContinueOnError)
	fs.StringVar(&listenPort, "listen", "localhost:36819", "IP and port to listen on")
	fs.StringVar(&storageBackend, "storage", "dory:", "DSN for storage backend")

	err := fs.Parse(os.Args[2:])
	if err != nil {
		return err
	}

	log.Printf("Starting server...")

	conf := plumbing.ServerConfig{
		Context:    context.Background(),
		StorageDSN: storageBackend,
	}
	s, err := plumbing.New(conf)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.Listen("tcp", listenPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s", listenPort)

	var srv http.Server
	srv.Handler = s
	return srv.Serve(ln)
}
