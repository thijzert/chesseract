package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/thijzert/chesseract/chesseract/client/httpclient"
	engine "github.com/thijzert/chesseract/internal/glengine"
	"github.com/thijzert/chesseract/internal/gui"
)

func glGame(conf *Config, args []string) error {
	fmt.Printf("Hello from version '%s'\n", engine.PackageVersion)

	// The actions didn't run because the Go compilation only runs when you change at least one Go file, you dummy.
	var autoquit int64

	rc := engine.DefaultConfig()
	rc.Logger = log.New(os.Stdout, "[life] ", log.Ltime|log.Lshortfile)

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

	// Draw a few circles
	red := color.RGBA{255, 0, 0, 255}
	yellow := color.RGBA{192, 192, 64, 255}
	for y := 0; y < debg.Pixels.Rect.Dy(); y++ {
		for x := 0; x < debg.Pixels.Rect.Dx(); x++ {
			xx, yy := x-120, y-120
			if xx*xx+yy*yy <= 10000 {
				debg.Pixels.SetRGBA(x, y, red)
			}
			xx, yy = x-1030, y-520
			if xx*xx+yy*yy <= 6400 {
				debg.Pixels.SetRGBA(x, y, yellow)
			}
		}
	}
	g, _ := os.Create("/tmp/debug-ui-out.png")
	png.Encode(g, debg.Pixels)

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
