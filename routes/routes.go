package routes

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tedbennett/battles/board"
	"github.com/tedbennett/battles/messages"
	"github.com/tedbennett/battles/relay"
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

func WebSocket(hub *relay.Relay, b *board.Board) echo.HandlerFunc {
	return func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()

			client := relay.NewClient(ws)
			defer func() {
				hub.Unregister <- client
			}()

			hub.Register <- client

			init := messages.NewInitMessage(b)
			initMsg := messages.NewMessage(init)

			websocket.Message.Send(ws, initMsg.Bytes())

			// Send broadcast messages, break on error
			for msg := range client.Messages {
				err := websocket.Message.Send(ws, msg)
				if err != nil {
					c.Logger().Error(err)
					break
				}
			}
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
