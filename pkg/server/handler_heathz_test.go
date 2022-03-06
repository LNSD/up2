package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"upload-presigned-url-provider/pkg/objectstore"
)

func TestHealthzHandlerObjectStoreConnected(t *testing.T) {
	/// Given
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	store := new(objectstore.MockObjectStore)
	store.On("Connect").Return(nil)

	e := echo.New()
	c := e.NewContext(req, rec)
	ctx := &CustomContext{c, store}
	server := Server{}

	/// When
	err := server.GetHealthz(ctx)

	/// Then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHealthzHandlerObjectStoreConnectionError(t *testing.T) {
	/// Given
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	store := new(objectstore.MockObjectStore)
	store.On("Connect").Return(fmt.Errorf("invalid credentials"))

	e := echo.New()
	c := e.NewContext(req, rec)
	ctx := &CustomContext{c, store}
	server := Server{}

	/// When
	err := server.GetHealthz(ctx)

	/// Then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}
