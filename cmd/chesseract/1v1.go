package main

import (
	"context"
	"fmt"
	"time"

	"github.com/thijzert/chesseract/chesseract"
	"github.com/thijzert/chesseract/chesseract/client"
	"github.com/thijzert/chesseract/chesseract/game"
	"github.com/thijzert/chesseract/internal/notimplemented"
)

type oneVoneServer struct {
	Game *game.Game
	W, B *oneVoneClient
}

func New1v1() *oneVoneServer {
	rv := &oneVoneServer{}
	rv.B = rv.NewClient()
	rv.B.colour = chesseract.BLACK

	rv.W = rv.NewClient()
	rv.W.colour = chesseract.WHITE

	return rv
}

func (o *oneVoneServer) NewClient() *oneVoneClient {
	rv := &oneVoneClient{
		server:   o,
		movesIn:  make(chan chesseract.Move, 1),
		resultIn: make(chan []float64, 1),
	}
	return rv
}

func (o *oneVoneServer) Run(ctx context.Context) {
	<-ctx.Done()
}

func (s *oneVoneServer) SubmitMove(c *oneVoneClient, _ *game.Game, m chesseract.Move) error {
	pat, ok := s.Game.Match.Board.At(m.From)
	if !ok {
		return client.ErrIllegalMove
	}

	if c == s.B {
		if pat.Colour != chesseract.BLACK || s.Game.Match.Board.Turn != chesseract.BLACK {
			return client.ErrNotYourTurn
		}
	} else if c == s.W {
		if pat.Colour != chesseract.WHITE || s.Game.Match.Board.Turn != chesseract.WHITE {
			return client.ErrNotYourTurn
		}
	} else {
		return client.ErrUnknownPlayer
	}

	newb, err := s.Game.Match.RuleSet.ApplyMove(s.Game.Match.Board, m)
	if err != nil {
		return client.ErrIllegalMove
	}

	// Reset the time from a centralised source
	m.Time = time.Since(s.Game.Match.StartTime)
	for _, m := range s.Game.Match.Moves {
		m.Time -= m.Time
	}

	s.Game.Match.Board = newb
	s.Game.Match.Moves = append(s.Game.Match.Moves, m)

	s.B.movesIn <- m
	s.W.movesIn <- m

	return nil
}

func (s *oneVoneServer) SubmitResult(c *oneVoneClient, _ *game.Game, result []float64) error {
	if s.Game.Result != nil {
		return client.ErrInvalidResult
	}
	if len(result) != 2 || result[0]+result[1] != 1.0 {
		return client.ErrInvalidResult
	}

	if c == s.B {
		if result[0] == 1.0 && result[1] == 0.0 {
			// Black forfeits
		} else {
			return fmt.Errorf("not implemented")
		}
	} else if c == s.W {
		if result[0] == 0.0 && result[1] == 1.0 {
			// White forfeits
		} else {
			// TODO: reach consensus on declaring a draw
			return fmt.Errorf("not implemented")
		}
	} else {
		return client.ErrUnknownPlayer
	}

	// Accept result
	s.Game.Result = append(s.Game.Result, result...)

	s.B.resultIn <- s.Game.Result
	s.W.resultIn <- s.Game.Result

	return nil
}

type oneVoneClient struct {
	server   *oneVoneServer
	colour   chesseract.Colour
	game     *game.Game
	movesIn  chan chesseract.Move
	resultIn chan []float64
}

// Me returns the object that represents the player at the server's end
func (o *oneVoneClient) Me() (game.Player, error) {
	// HACK: since this is just a 1v1, the only two that currently exist also happen to be in-game.
	for _, pl := range o.server.Game.Players {
		if pl.PlayingAs == o.colour {
			return pl.Player, nil
		}
	}

	return game.Player{}, fmt.Errorf("player does not exist")
}

// AvailablePlayers returns the list of players available for a match
func (o *oneVoneClient) AvailablePlayers(context.Context) ([]game.Player, error) {
	// HACK: since this is just a 1v1, the only two that currently exist also happen to be in-game.
	rv := make([]game.Player, 0, 1)
	for _, pl := range o.server.Game.Players {
		if pl.PlayingAs != o.colour {
			rv = append(rv, pl.Player)
		}
	}

	return rv, nil
}

func (o *oneVoneClient) ActiveGames(context.Context) ([]client.GameSession, error) {
	if o.server.Game == nil {
		return nil, nil
	} else if o.game == nil {
		o.game = &game.Game{}
		o.game.Players = append(o.game.Players, o.server.Game.Players...)
		o.game.Match.RuleSet = o.server.Game.Match.RuleSet
		o.game.Match.Board.Turn = o.server.Game.Match.Board.Turn
		o.game.Match.Board.Pieces = append(o.game.Match.Board.Pieces, o.server.Game.Match.Board.Pieces...)
		o.game.Match.StartTime = o.server.Game.Match.StartTime
		o.game.Match.Moves = append(o.game.Match.Moves, o.server.Game.Match.Moves...)
	}

	return []client.GameSession{o}, nil
}

func (o *oneVoneClient) NewGame(ctx context.Context, rs chesseract.RuleSet, players []game.Player) (client.GameSession, error) {
	if o.server.Game == nil {
		o.server.Game = &game.Game{
			Match: chesseract.Match{
				RuleSet:   rs,
				Board:     rs.DefaultBoard(),
				StartTime: time.Now(),
			},
		}
		for i, c := range o.server.Game.Match.RuleSet.PlayerColours() {
			if len(players) > i {
				o.server.Game.Players = append(o.server.Game.Players, game.MatchPlayer{
					Player:    players[i],
					PlayingAs: c,
				})
			}
		}
	}

	// Create client-side game copy
	o.ActiveGames(ctx)

	return o, nil
}

// Game returns the Game object of this session
func (o *oneVoneClient) Game() *game.Game {
	return o.game
}

func (o *oneVoneClient) PlayingAs() chesseract.Colour {
	return o.colour
}

// SubmitMove submits a move by this player.
func (o *oneVoneClient) SubmitMove(ctx context.Context, move chesseract.Move) error {
	return o.server.SubmitMove(o, o.game, move)
}

// NextMove waits until a move occurs, and returns it. This comprises moves
// made by all players, not just opponents.
func (o *oneVoneClient) NextMove(ctx context.Context) (chesseract.Move, error) {
	select {
	case <-ctx.Done():
		return chesseract.Move{}, ctx.Err()

	case m := <-o.movesIn:
		newb, err := o.game.Match.RuleSet.ApplyMove(o.game.Match.Board, m)
		if err != nil {
			return chesseract.Move{}, client.ErrIllegalMove
		}

		o.game.Match.Board = newb
		o.game.Match.Moves = append(o.game.Match.Moves, m)
		return m, nil
	}
}

// ProposeResult submits a possible final outcome for this game, which all
// opponents can evaluate and accept or reject. One can accept a proposed
// result by proposing the same result again.
// Proposing a nil or zero result is construed as rejecting a proposition.
func (o *oneVoneClient) ProposeResult(ctx context.Context, result []float64) error {
	return o.server.SubmitResult(o, o.game, result)
}

// NextProposition waits until a result is proposed, and returns it.
func (o *oneVoneClient) NextProposition(ctx context.Context) ([]float64, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case res := <-o.resultIn:
		return res, nil
	}
}

// GetResult retrieves the result for this game
func (o *oneVoneClient) GetResult(context.Context) ([]float64, error) {
	return nil, notimplemented.Error()
}
