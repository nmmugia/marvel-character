package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nmmugia/marvel-character/models"
)

// Message function is to build response per standard
func Message(errx models.Errorx, data interface{}) (res map[string]interface{}) {
	res = map[string]interface{}{
		"status":  errx.Status,
		"message": errx.Message,
	}
	if data != nil {
		res["data"] = data
	}
	if errx.Err != nil {
		res["internalMessage"] = errx.Err.Error()
	}
	return res
}

// Response function is to encode response
func Response(w http.ResponseWriter, data interface{}, errx models.Errorx) {
	w.Header().Add("Content-Type", "application/json")
	if errx.Err != nil {
		json.NewEncoder(w).Encode(Message(errx, data))
	} else if errx.Status != 0 && errx.Message != "" {
		json.NewEncoder(w).Encode(Message(errx, nil))
	} else if data == nil && errx.Status == 0 {
		json.NewEncoder(w).Encode(Message(models.Errorx{
			Message: "No Content",
			Status:  http.StatusNoContent,
		}, nil))
	} else {
		json.NewEncoder(w).Encode(Message(models.Errorx{
			Message: "Success",
			Status:  http.StatusOK,
		}, data))
	}
}

// StringToMD5 function is to convert/hash string to MD5 format
func StringToMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// StringToInt to casting string to integer, first param would be value, and the second one is default returned value
func StringToInt(value string, def int64) int64 {
	r, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return def
	}
	return r
}

// ParseFromString to parse time from string
func ParseFromString(value string) (res time.Time, err error) {
	for _, f := range []string{
		"02-01-2006",
		"2006-01-02",
		"20060102150405",
		"20060102",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05 Z07:00",
		"2006-01-02Z07:00",
		"20060102150405Z07:00",
		"20060102Z07:00",
		"2006-01-02 15:04",
		"2006-01-02 15",
		"2006-01",
		"2006",
	} {
		res, err = time.ParseInLocation(f, value, time.Local)
		if err == nil {
			return res, err
		}
	}

	if IsInteger(value) {
		res = time.Unix(StringToInt(value, 0), 0)
	}

	return res, err
}

// IsInteger : Is Integer ?
func IsInteger(v string) bool {
	if v == "" {
		return false
	}

	a := "1234567890"
	for _, v := range v {
		if !strings.Contains(a, string(v)) {
			return false
		}
	}

	return true
}
