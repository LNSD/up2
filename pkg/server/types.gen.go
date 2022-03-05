// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package server

// DownloadRequestBody defines model for DownloadRequestBody.
type DownloadRequestBody struct {
	// URL expiration time in seconds
	Expiration *int `json:"expiration,omitempty"`

	// Identifier of the file to download
	Id string `json:"id"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	// The error message
	Message string `json:"message"`
}

// PreSignedURL defines model for PreSignedURL.
type PreSignedURL struct {
	// URL expiration time in seconds
	Expiration int `json:"expiration"`

	// URL identifier
	Id string `json:"id"`

	// The presigned URL
	URL string `json:"url"`
}

// UploadRequestBody defines model for UploadRequestBody.
type UploadRequestBody struct {
	// URL expiration time in seconds
	Expiration *int `json:"expiration,omitempty"`
}

// PostDownloadJSONBody defines parameters for PostDownload.
type PostDownloadJSONBody DownloadRequestBody

// PostUploadJSONBody defines parameters for PostUpload.
type PostUploadJSONBody UploadRequestBody

// PostDownloadJSONRequestBody defines body for PostDownload for application/json ContentType.
type PostDownloadJSONRequestBody PostDownloadJSONBody

// PostUploadJSONRequestBody defines body for PostUpload for application/json ContentType.
type PostUploadJSONRequestBody PostUploadJSONBody
