package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/thijzert/chesseract/chesseract"
	"github.com/thijzert/chesseract/chesseract/client"
	"github.com/thijzert/chesseract/chesseract/client/httpclient"
	"github.com/thijzert/chesseract/chesseract/game"
)

var consoleMutex sync.Mutex

func consoleGame(conf *Config, args []string) error {
	logVerbose := false
	clientConf := httpclient.ClientConfig{}

	consoleSettings := flag.NewFlagSet("consoleClient", flag.ContinueOnError)
	consoleSettings.StringVar(&clientConf.ServerURI, "server", "", "URI to multiplayer server")
	consoleSettings.StringVar(&clientConf.Username, "username", "", "Online username")
	consoleSettings.BoolVar(&logVerbose, "v", false, "Verbosely log all requests")
	err := consoleSettings.Parse(args)
	if err != nil {
		return err
	}

	if logVerbose {
		clientConf.VerboseRequestLogging = os.Stdout
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var c client.Client
	c, err = httpclient.New(ctx, clientConf)
	if err != nil {
		return err
	}

	g, err := c.NewGame(ctx, []game.Player{
		{Name: "white"},
		{Name: "black"},
	})
	if err != nil {
		return err
	}
	cc := newConsoleClient(c, chesseract.BLACK)

	return cc.Run(ctx, g)
}

func consoleLocalMultiplayer(conf *Config, args []string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	players := []game.Player{
		{Name: "white"},
		{Name: "black"},
	}

	errs := make(chan error, 4)

	s := New1v1()

	var run = func(c client.Client, col chesseract.Colour) {
		g, err := c.NewGame(ctx, players)
		if err != nil {
			errs <- err
			return
		}

		cc := newConsoleClient(c, col)
		errs <- cc.Run(ctx, g)
	}

	go run(s.B, chesseract.BLACK)
	go run(s.W, chesseract.WHITE)

	var rv error

	for i := 0; i < 2; i++ {
		err := <-errs
		cancel()
		if rv == nil && err != nil {
			rv = err
		}
	}
	cancel()
	return rv
}

type consoleClient struct {
	Client client.Client
	PlayAs chesseract.Colour
}

func newConsoleClient(c client.Client, playAs chesseract.Colour) consoleClient {
	return consoleClient{
		Client: c,
		PlayAs: playAs,
	}
}

func (cc consoleClient) Run(ctx context.Context, g *game.Game) error {
	for ctx.Err() == nil {
		for g.Match.Board.Turn != cc.PlayAs {
			_, err := cc.Client.NextMove(ctx, g)
			if err != nil {
				return err
			}
		}

		consoleMutex.Lock()
		g.Match.DebugDump(os.Stdout, nil)
		consoleMutex.Unlock()

		var move chesseract.Move

		for {
			consoleMutex.Lock()
			fmt.Printf("Enter move for %6s: ", cc.PlayAs)

			var sFrom, sTo string
			n, _ := fmt.Scanf("%s %s\n", &sFrom, &sTo)
			consoleMutex.Unlock()
			if n == 0 {
				continue
			}
			if n == 1 {
				if sFrom == "forfeit" || sFrom == "quit" {
					return fmt.Errorf("forfeiting is not implemented")
				}
			}

			from, err := g.Match.RuleSet.ParsePosition(sFrom)
			if err != nil {
				fmt.Printf("error parsing '%s': %v\n", sFrom, err)
				continue
			}
			piece, _ := g.Match.Board.At(from)
			to, err := g.Match.RuleSet.ParsePosition(sTo)
			if err != nil {
				fmt.Printf("error parsing '%s': %v\n", sTo, err)
				continue
			}

			move = chesseract.Move{
				PieceType: piece.PieceType,
				From:      from,
				To:        to,
			}
			_, err = g.Match.RuleSet.ApplyMove(g.Match.Board, move)
			if err != nil {
				fmt.Printf("applying move '%s'-'%s': %v\n", sFrom, sTo, err)
				continue
			}

			break
		}

		cc.Client.SubmitMove(ctx, g, move)
		type moveErr struct {
			Move chesseract.Move
			Err  error
		}
		ch := make(chan moveErr)
		go func() {
			otherMove, err := cc.Client.NextMove(ctx, g)
			ch <- moveErr{otherMove, err}
			close(ch)
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case mv := <-ch:
			if mv.Err != nil {
				return mv.Err
				// TODO: Maybe the server just thinks this is illegal, and we should keep trying?
			}
			if !mv.Move.From.Equals(move.From) || !mv.Move.To.Equals(move.To) {
				return client.ErrShenanigans
			}
		}
	}
	return ctx.Err()
}
