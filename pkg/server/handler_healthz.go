package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s Server) GetHealthz(c echo.Context) error {
	ctx := c.(*CustomContext)
	objectStore := ctx.ObjectStore

	if err := objectStore.Connect(); err != nil {
		return ctx.JSON(
			http.StatusServiceUnavailable,
			ErrorResponse{"object store connection not available"},
		)
	}

	return ctx.NoContent(http.StatusOK)
}
