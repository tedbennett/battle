package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"

	"github.com/tedbennett/battles/board"
	"github.com/tedbennett/battles/messages"
	"github.com/tedbennett/battles/relay"
	"github.com/tedbennett/battles/routes"
	"github.com/tedbennett/battles/templates"
	"github.com/tedbennett/battles/utils"
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

	b := board.NewBarBoard(20,
		map[int]utils.Color{
			0: utils.ColorFromString(board.Team1Color),
			1: utils.ColorFromString(board.Team2Color),
		},
	)
	channel := make(chan []byte)

	hub := relay.NewRelay()

	go hub.Run()

	go func() {
		defer close(channel)
		for {
			time.Sleep(time.Millisecond * 500)
			diffs := b.TickBar()
			e.Logger.Info("Sending board message")
			msg := messages.NewMessage(messages.NewPartialMessage(diffs))
			hub.Broadcast <- msg.Bytes()
		}
	}()
	e.GET("/", routes.Home(&b))
	e.GET("/ws", routes.WebSocket(hub, &b))
	e.Logger.Fatal(e.Start(":8000"))
}
