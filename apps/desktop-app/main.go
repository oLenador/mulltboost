package main

import (
	"context"
	"embed"
	"log"

	"github.com/oLenador/mulltbost/internal/app/container"
	"github.com/oLenador/mulltbost/internal/app/handlers"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	container := container.NewContainer()

	if err := container.RegisterPlugins(); err != nil {
		log.Fatal("Failed to register plugins:", err)
	}

	// Cria handlers
	optimizationHandler := handlers.NewOptimizationHandler(container)
	monitoringHandler := handlers.NewMonitoringHandler(container)
	systemHandler := handlers.NewSystemHandler(container)

	err := wails.Run(&options.App{
		Title:     "Mulltboost",
		Width:     1024,
		Height:    768,
		MinWidth:  1280,
		MinHeight: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
        OnStartup: func(ctx context.Context) {
            optimizationHandler.SetContext(ctx)
            monitoringHandler.SetContext(ctx)
            systemHandler.SetContext(ctx)
        },
        Bind: []interface{}{
            optimizationHandler,
            monitoringHandler,
            systemHandler,
        },
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
