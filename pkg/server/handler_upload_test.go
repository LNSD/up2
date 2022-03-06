package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jaswdr/faker"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"upload-presigned-url-provider/pkg/objectstore"
)

func fakeUploadPreSignedURL(expiration time.Duration) objectstore.PreSignedURL {
	faker := faker.New()
	return objectstore.PreSignedURL{
		URL:        faker.Internet().URL(),
		Expiration: expiration,
	}
}

func TestUploadHandlerObjectStoreError(t *testing.T) {
	/// Given
	store := new(objectstore.MockObjectStore)
	store.On("GetUploadURL", mock.Anything).Return(nil, fmt.Errorf("object store failure"))

	// HTTP request
	reqParams := UploadRequestBody{}
	reqBody, _ := json.Marshal(reqParams)
	req := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Server
	e := echo.New()
	c := e.NewContext(req, rec)
	ctx := &CustomContext{c, store}
	server := Server{}

	/// When
	err := server.PostUpload(ctx)

	/// Then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "object store cannot generate a pre-signed url", resp.Message)
}

func TestUploadHandlerNoExpirationTime(t *testing.T) {
	/// Given
	defaultExpiration := 10
	preSignedURL := fakeUploadPreSignedURL(time.Duration(defaultExpiration) * time.Second)
	store := new(objectstore.MockObjectStore)
	store.On("GetUploadURL", mock.Anything).Return(&preSignedURL, nil)

	// HTTP request
	reqParams := UploadRequestBody{}
	reqBody, _ := json.Marshal(reqParams)
	req := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Server
	e := echo.New()
	c := e.NewContext(req, rec)
	ctx := &CustomContext{c, store}
	server := Server{}

	/// When
	err := server.PostUpload(ctx)

	/// Then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp PreSignedURL
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, defaultExpiration, resp.Expiration)
	assert.Equal(t, preSignedURL.URL, resp.URL)
	assert.NotEmpty(t, resp.Id)
}

func TestUploadHandlerExplicitExpirationTime(t *testing.T) {
	/// Given
	expiration := 10

	// Object store
	preSignedURL := fakeUploadPreSignedURL(time.Duration(expiration) * time.Second)
	store := new(objectstore.MockObjectStore)
	store.On("GetUploadURLWithExpiration", mock.Anything, mock.Anything).Return(&preSignedURL, nil)

	// HTTP request
	reqParams := UploadRequestBody{&expiration}
	reqBody, _ := json.Marshal(reqParams)
	req := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Server
	e := echo.New()
	c := e.NewContext(req, rec)
	ctx := &CustomContext{c, store}
	server := Server{}

	/// When
	err := server.PostUpload(ctx)

	/// Then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp PreSignedURL
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, expiration, resp.Expiration)
	assert.Equal(t, preSignedURL.URL, resp.URL)
	assert.NotEmpty(t, resp.Id)
}

func TestUploadHandlerInvalidExpirationTime(t *testing.T) {
	/// Given
	expiration := -1

	// Object store
	store := new(objectstore.MockObjectStore)

	// HTTP request
	reqParams := UploadRequestBody{&expiration}
	reqBody, _ := json.Marshal(reqParams)
	req := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Server
	e := echo.New()
	c := e.NewContext(req, rec)
	ctx := &CustomContext{c, store}
	server := Server{}

	/// When
	err := server.PostUpload(ctx)

	/// Then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "expiration parameter must be an integer bigger than 0", resp.Message)
}

func TestUploadHandlerInvalidRequestBody(t *testing.T) {
	/// Given

	// Object store
	store := new(objectstore.MockObjectStore)

	// HTTP request
	body := bytes.NewBufferString(`{`)
	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Server
	e := echo.New()
	c := e.NewContext(req, rec)
	ctx := &CustomContext{c, store}
	server := Server{}

	/// When
	err := server.PostUpload(ctx)

	/// Then
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp ErrorResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "invalid request body", resp.Message)
}
