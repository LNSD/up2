package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s Server) GetHealthz(c echo.Context) error {
	ctx := c.(*CustomContext)

	// Check if Object Store connection is ready
	if ctx.ObjectStore == nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			ErrorResponse{"object store connection not available"},
		)
	}

	return ctx.NoContent(http.StatusOK)
}
