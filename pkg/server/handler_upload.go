package server

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	"upload-presigned-url-provider/pkg/objectstore"
)

func (s Server) PostUpload(c echo.Context) error {
	ctx := c.(*CustomContext)

	// Check if Object Store connection is ready
	if ctx.ObjectStore == nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			ErrorResponse{"object store connection not available"},
		)
	}
	objectStore := *ctx.ObjectStore

	// Get URL unique identifier
	objectId, err := uuid.NewRandom()
	if err != nil {
		c.Logger().Error(err)
		return ctx.JSON(
			http.StatusInternalServerError,
			ErrorResponse{"cannot generated a unique url id"},
		)
	}

	// Get URL expiration time
	var req *UploadRequestBody
	if err := c.Bind(req); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			ErrorResponse{"invalid request body"},
		)
	}

	if req.Expiration != nil {
		expiration := *req.Expiration
		if expiration <= 0 {
			return ctx.JSON(
				http.StatusBadRequest,
				ErrorResponse{"expiration parameter must be an integer bigger than 0"},
			)
		}
	}

	// Generate a pre-signed upload URL
	objectKey := objectId.String()

	var preSignedURL *objectstore.PreSignedUrl
	if req.Expiration == nil {
		preSignedURL, err = objectStore.GetUploadUrl(objectKey)
	} else {
		expiration := time.Duration(*req.Expiration) * time.Second
		preSignedURL, err = objectStore.GetUploadUrlWithExpiration(objectKey, expiration)
	}

	if err != nil {
		c.Logger().Error(err)
		return ctx.JSON(
			http.StatusInternalServerError,
			ErrorResponse{"object store cannot generate a pre-signed url"},
		)
	}

	return ctx.JSON(http.StatusOK, PreSignedUrl{
		Url:        preSignedURL.Url,
		Expiration: int(preSignedURL.Expiration.Seconds()),
		Id:         objectKey,
	})
}
