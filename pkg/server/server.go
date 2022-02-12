package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"upload-presigned-url-provider/pkg/objectstore"
)

type Server struct {
	*Config
	*echo.Echo
}

func New(store *objectstore.ObjectStore, config *Config) Server {
	e := echo.New()
	e.HideBanner = true
	e.Use(NewCustomContextMiddleware(store))
	RegisterHandlers(e, Server{})
	return Server{config, e}
}

func (s Server) Start() {
	err := s.Echo.Start(s.Config.Address())
	if err != nil && err != http.ErrServerClosed {
		// TODO: Review this. Maybe, return this error and handle it properly
		s.Echo.Logger.Fatal("shutting down the Server")
	}
}
func (s Server) Shutdown(ctx context.Context) {
	if err := s.Echo.Shutdown(ctx); err != nil {
		// TODO: Review this. Maybe, return this error and handle it properly
		s.Echo.Logger.Fatal(err)
	}
}
