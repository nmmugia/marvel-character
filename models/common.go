package models

import (
	"errors"
	"net/http"

	"github.com/jmoiron/sqlx/types"
)

// var db *sqlx.DB

type (
	MarvelsResult struct {
		Code    int16          `json:"code"`
		Data    types.JSONText `json:"data"`
		Message string         `json:"message"`
	}
	Errorx struct {
		Err     error
		Status  int16
		Message string
	}
)

const (
	MessageInternalServerError = "Kesalahan terjadi di server|Internal server error"
	MessageBadRequest          = "Mohon cek ulang request anda|Kindly check you request"
	MessageNotFound            = "Yang anda cari tidak ditemukan|Not Found"
	MessageForbidden           = "Akses Dilarang|Forbidden access"
	MessageStatusUnauthorized  = "Akses Anda tidak terautorisasi|Unauthorized Access"
)

// CreateErrx is a func to create errorx
func CreateErrx(status int16, err error, message string) Errorx {
	if status <= 0 || status >= 600 {
		status = http.StatusInternalServerError
	}
	if err == nil {
		err = errors.New("internal server error")
	}

	// if message is not set
	if len(message) <= 0 {
		message = "Internal Server Error"
	}
	return Errorx{
		Status:  status,
		Err:     err,
		Message: message,
	}
}
