package server

import (
	"github.com/labstack/echo/v4"
	"upload-presigned-url-provider/pkg/objectstore"
)

type CustomContext struct {
	echo.Context
	ObjectStore objectstore.ObjectStore
}

func NewCustomContextMiddleware(objectStore objectstore.ObjectStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			cc := &CustomContext{
				Context:     ctx,
				ObjectStore: objectStore,
			}
			return next(cc)
		}
	}
}
