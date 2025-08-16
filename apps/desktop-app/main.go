package main

import (
	"context"
	"embed"

	"github.com/oLenador/mulltbost/internal/app/container"
	"github.com/oLenador/mulltbost/internal/app/handlers"
	"github.com/oLenador/mulltbost/internal/config"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

const AppName = "mulltboost"

func main() {

	container, err := container.NewContainer()
	
	if err != nil {
		println("Error:", err.Error())
	}

	// Cria handlers
	boosterHandler := handlers.NewBoosterHandler(container)
	monitoringHandler := handlers.NewMonitoringHandler(container)
	systemHandler := handlers.NewSystemHandler(container)

	err = wails.Run(&options.App{
		Title:     config.AppTitle,
		Width:     1024,
		Height:    628,
		MinWidth:  1050,
		MinHeight: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
        OnStartup: func(ctx context.Context) {
            boosterHandler.SetContext(ctx)
            monitoringHandler.SetContext(ctx)
            systemHandler.SetContext(ctx)
        },
        Bind: []interface{}{
            boosterHandler,
            monitoringHandler,
            systemHandler,
        },
		Windows: &windows.Options{
			DisableFramelessWindowDecorations: true,
			DisableWindowIcon: true,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
