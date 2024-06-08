package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type HomeTemplateData struct {
	Colors [][]string
}

func Home() echo.HandlerFunc {
	return func(c echo.Context) error {
		// colors := b.Colors()
		// data := HomeTemplateData{colors}
		return c.Render(http.StatusOK, "index", struct{}{})
	}
}

func WebSocket(channel <-chan []byte, init []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()

			// Send init message
			err := websocket.Message.Send(ws, init)
			if err != nil {
				c.Logger().Error(err)
			}

			// On tick, send board update
			for msg := range channel {
				err := websocket.Message.Send(ws, msg)
				if err != nil {
					c.Logger().Error(err)
				}
			}
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
