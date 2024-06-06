package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tedbennett/battles/board"
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
