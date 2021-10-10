package api

import (
	"log"
	"proxy/pkg/database"

	"github.com/aerogo/aero"
)

func configure(app *aero.Application, handlers *Handler) *aero.Application {
	app.Get("/api/requests", handlers.RequestAll)
	app.Get("/api/requests/:id", handlers.RequestOne)
	app.Get("/api/requests/:id/repeat", handlers.Repeat)
	app.Get("/api/requests/:id/scan", handlers.Scan)
	return app
}

func Run() {
	db := database.GetTarantool()
	handlers := &Handler{db}

	app := aero.New()
	app.OnEnd(func() {
		database.CloseTarantool(db)
	})
	app = configure(app, handlers)
	log.Println(app.Config)
	app.Run()
}
