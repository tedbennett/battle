package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tedbennett/battles/board"
	"golang.org/x/net/websocket"
)

type HomeTemplateData struct {
	Colors [][]string
}

func Home(b *board.Board) echo.HandlerFunc {
	return func(c echo.Context) error {
		colors := b.Colors()
		data := HomeTemplateData{colors}
		return c.Render(http.StatusOK, "index", data)
	}
}
func WebSocket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			// Write
			err := websocket.Message.Send(ws, "Hello, Client!")
			if err != nil {
				c.Logger().Error(err)
			}

			// Read
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
