package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thijzert/chesseract/chesseract"
	"github.com/thijzert/chesseract/chesseract/client"
	"github.com/thijzert/chesseract/chesseract/client/httpclient"
	"github.com/thijzert/chesseract/chesseract/game"
	engine "github.com/thijzert/chesseract/internal/glengine"
	"github.com/thijzert/chesseract/internal/gui"
)

func glGame(conf *Config, args []string) error {
	fmt.Printf("Hello from version '%s'\n", engine.PackageVersion)

	var autoquit int64

	rc := engine.DefaultConfig()
	rc.Logger = log.New(os.Stdout, "[gl] ", log.Ltime|log.Lshortfile)

	clientConf := httpclient.ClientConfig{}

	glSettings := flag.NewFlagSet("glClient", flag.ContinueOnError)
	glSettings.StringVar(&clientConf.ServerURI, "server", "", "URI to multiplayer server")
	glSettings.StringVar(&clientConf.Username, "username", "", "Online username")
	glSettings.IntVar(&rc.WindowWidth, "w", 1280, "Window width")
	glSettings.IntVar(&rc.WindowHeight, "h", 720, "Window height")
	glSettings.Int64Var(&autoquit, "autoquit", 0, "Automatically quit after x ms (0 to disable)")
	err := glSettings.Parse(args)
	if err != nil {
		return err
	}

	ctx := context.Background()

	if autoquit > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(autoquit)*time.Millisecond)
		defer cancel()
	}
	eng := rc.NewEngine(ctx)

	debg := gui.GUIContext{
		Pixels: image.NewRGBA(image.Rect(0, 0, rc.WindowWidth, rc.WindowHeight)),
	}
	eng.GUI.AddLayer(GUI_DEBUG, debg)

	go func() {
		er := func() error {
			var c client.Client
			c, err = httpclient.New(ctx, clientConf)
			if err != nil {
				return err
			}

			var g client.GameSession
			ag, err := c.ActiveGames(ctx)
			if err != nil {
				return err
			}

			if len(ag) > 0 {
				g = ag[0]
			} else {
				g, err = c.NewGame(ctx, []game.Player{
					{Name: "alice"},
					{Name: "bob"},
				})
				if err != nil {
					return err
				}
			}

			cc := glClient{
				Session:         g,
				RenderingEngine: eng,
			}

			return cc.Run(ctx)
		}()
		if er != nil {
			log.Print(er)
		}
	}()

	// FIXME: The shutdown needs to happen in the same OS thread as the engine itself.
	//        Better to just defer the cleanup stage inside Run(). As a bonus, less nasty error wrangling below.
	err = eng.Run()
	er1 := eng.Shutdown()
	if err != nil {
		return err
	}
	return er1
}

const (
	GUI_DEBUG engine.LayerName = iota
	GUI_HUD
	GUI_MENU
	GUI_ALERT
)

type glClient struct {
	Session         client.GameSession
	RenderingEngine *engine.Engine
}

func (cc glClient) Run(ctx context.Context) error {
	playingAs := cc.Session.PlayingAs()
	g := cc.Session.Game()
	for ctx.Err() == nil {
		if g.Match.Board.Turn != playingAs {
			fmt.Printf("Waiting for opponent\n")
		}
		for g.Match.Board.Turn != playingAs {
			_, err := cc.Session.NextMove(ctx)
			if err != nil {
				return err
			}
		}

		cc.RenderBoard()

		var move chesseract.Move

		for {
			fmt.Printf("Enter move for %6s: ", playingAs)

			var sFrom, sTo string
			n, _ := fmt.Scanf("%s %s\n", &sFrom, &sTo)
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

		err := cc.Session.SubmitMove(ctx, move)
		if err != nil {
			return err
		}

		type moveErr struct {
			Move chesseract.Move
			Err  error
		}
		ch := make(chan moveErr)
		go func() {
			otherMove, err := cc.Session.NextMove(ctx)
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

			cc.RenderBoard()
		}
	}
	return ctx.Err()
}

func (cc glClient) RenderBoard() {
	cc.RenderingEngine.ClearEntities()
	rs := cc.Session.Game().Match.RuleSet

	for i, piece := range cc.Session.Game().Match.Board.Pieces {
		name, tileIdx, ok := resolveModel(piece)
		if mr, isMR := rs.(ModelResolver); isMR {
			name, tileIdx, ok = mr.ResolveModel(piece)
		}
		if !ok {
			continue
		}
		position := mgl32.Vec3{float32(i), -2, -6}

		if pos, ok := piece.Position.(Positioner); ok {
			x, y, z := pos.WorldPosition()
			position = mgl32.Vec3{x, y, z}
		}

		cc.RenderingEngine.Entities = append(cc.RenderingEngine.Entities, engine.Entity{
			ModelName: name,
			Position:  position,
			Scale:     mgl32.Vec3{1, 1, 1},
			TileIndex: tileIdx,
		})
	}

	cc.RenderingEngine.SwapEntities()
}

func resolveModel(piece chesseract.Piece) (modelName string, tileIdx int, ok bool) {
	modelNames := map[chesseract.PieceType]string{
		chesseract.PAWN:   "pawn",
		chesseract.ROOK:   "rook",
		chesseract.KNIGHT: "knight",
		chesseract.BISHOP: "bishop",
		chesseract.QUEEN:  "queen",
		chesseract.KING:   "king",
	}

	idxMap := map[chesseract.Colour]int{
		chesseract.WHITE: 0,
		chesseract.BLACK: 1,
	}

	modelName, ok = modelNames[piece.PieceType]
	tileIdx = idxMap[piece.Colour]
	return
}

type Positioner interface {
	WorldPosition() (float32, float32, float32)
}

type ModelResolver interface {
	ResolveModel(piece chesseract.Piece) (modelName string, tileIdx int, ok bool)
}
