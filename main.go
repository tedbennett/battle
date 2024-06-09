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

	b := board.NewBarBoard(20)
	colors := map[int]utils.Color{
		0: utils.ColorFromString(board.Team1Color),
		1: utils.ColorFromString(board.Team2Color),
	}

	init := messages.NewInitMessage(colors, &b)
	initMsg := messages.NewMessage(20, init)
	initMsg.Debug()
	channel := make(chan []byte)
	go func(channel chan<- []byte) {
		defer close(channel)
		for {
			time.Sleep(time.Millisecond * 500)
			diffs := b.TickBar()
			e.Logger.Info("Sending board message")
			msg := messages.NewMessage(255, messages.NewPartialMessage(diffs))
			channel <- msg.Bytes()
		}
	}(channel)
	e.GET("/", routes.Home(&b))
	e.GET("/ws", routes.WebSocket(channel, initMsg.Bytes()))
	e.Logger.Fatal(e.Start(":8000"))
}
