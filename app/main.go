package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"reqx/app/wailsapp"
)

// all: forces embedding of dot-prefixed files too, so a checked-in
// placeholder under frontend/dist survives before the first `npm run build`
// regenerates the real bundle. See frontend/dist/index.html.
//
//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := wailsapp.NewApp()

	err := wails.Run(&options.App{
		Title:  "ReqX",
		Width:  1200,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Bind: app.BindTargets(),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "reqx-app:", err)
		os.Exit(1)
	}
}
