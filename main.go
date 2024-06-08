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

	b := board.NewBoard(5, 0)
	colors := map[int]string{0: board.Team1Color, 1: board.Team2Color}
	init := messages.NewInitMessage(colors, b.Squares)
	initMsg, _ := json.Marshal(init)
	channel := make(chan []byte)
	go func(channel chan<- []byte) {
		defer close(channel)
		count := 0
		for {
			time.Sleep(time.Millisecond * 2000)
			e.Logger.Info("Sending board message")
			b := board.NewBoard(5, 0)
			b.Squares[count/5][count%5] = 1
			msg := messages.NewBoardMessage(b.Squares)
			json, _ := json.Marshal(msg)
			channel <- json
			if count < 24 {
				count += 1
			} else {
				count = 0
			}
		}
	}(channel)
	e.GET("/", routes.Home())
	e.GET("/ws", routes.WebSocket(channel, initMsg))
	e.Logger.Fatal(e.Start(":8000"))
}
