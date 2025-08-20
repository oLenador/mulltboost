package main

import (
	"embed"
	"log"

	"github.com/oLenador/mulltbost/internal/app/container"
	"github.com/oLenador/mulltbost/internal/app/handlers"
	"github.com/oLenador/mulltbost/internal/config"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
//go:embed all:frontend/dist
var assets embed.FS

const AppName = "mulltboost"

func main() {

	app := application.New(application.Options{
		Name:        config.AppTitle,
		Description: "mulltboost desktop application",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Linux: application.LinuxOptions{
			ProgramName:                  config.AppTitle,
			DisableQuitOnLastWindowClosed: true,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Now initialize your DI/container and pass the app so the container can interact with it.
	svcContainer, err := container.NewContainer(app)
	if err != nil {
		log.Fatal("Error initializing container:", err)
	}

	metricsHandler := handlers.NewMetricsHandler(svcContainer)
	boosterHandler := handlers.NewBoosterHandler(svcContainer)
	systemHandler := handlers.NewSystemHandler(svcContainer)

	app.RegisterService(application.NewService(metricsHandler))
	app.RegisterService(application.NewService(boosterHandler))
	app.RegisterService(application.NewService(systemHandler))


	// Create the main window with the necessary options
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  config.AppTitle,
		Width:  1024,
		Height: 628,
		Windows: application.WindowsWindow{
			DisableFramelessWindowDecorations: true,
			DisableIcon:                 true,
			BackdropType: application.Acrylic,
		},
		MinWidth:        1050,
		MinHeight:       768,
		BackgroundColour: application.NewRGB(0, 0, 0), 
		URL:             "/",
	})

	// Run the application. This blocks until the application has been exited.
	err = app.Run()
	if err != nil {
		log.Fatal("Error running application:", err)
	}
}