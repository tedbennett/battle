package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"

	"github.com/tedbennett/battles/board"
	"github.com/tedbennett/battles/messages"
	"github.com/tedbennett/battles/routes"
	"github.com/tedbennett/battles/templates"
)

//go:generate npm run build
func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	e.Use(slogecho.New(logger))

	e.Static("static/css", "static/css")
	e.Static("static/js", "static/js")
	templates.NewTemplateRenderer(e, "static/*.html")

	b := board.NewConwayBoard(200)
	colors := map[int]string{0: board.Team1Color, 1: board.Team2Color}
	init := messages.NewInitMessage(colors, b.Squares)
	initMsg, _ := json.Marshal(init)
	channel := make(chan []byte)
	go func(channel chan<- []byte) {
		defer close(channel)
		for {
			time.Sleep(time.Millisecond * 500)
			b.Tick()
			e.Logger.Info("Sending board message")
			msg := messages.NewBoardMessage(b.Squares)
			json, _ := msg.Marshal()
			channel <- json
		}
	}(channel)
	e.GET("/", routes.Home())
	e.GET("/ws", routes.WebSocket(channel, initMsg))
	e.Logger.Fatal(e.Start(":8000"))
}
