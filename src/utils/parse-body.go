package utils

import (
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func ParseBody[T any](r *http.Request) (T, error) {
	body, _ := io.ReadAll(r.Body)
	var data T
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, &data)
	return data, err
}
