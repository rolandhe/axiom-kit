package zinc_sdk

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var ErrDocumentNoExist = errors.New("document no exist")
var ErrIdNotFound = errors.New("id not found")

func idNotFoundRespCallback(resp *http.Response) (bool, error) {
	if resp.StatusCode != http.StatusBadRequest {
		return false, nil
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return true, err
	}
	if len(content) == 0 {
		return true, errors.New(resp.Status)
	}

	body := struct {
		Error string `json:"error"`
	}{}
	err = json.Unmarshal(content, &body)
	if err != nil {
		return true, err
	}
	if body.Error == "id not found" {
		return true, ErrIdNotFound
	}
	return true, errors.New(resp.Status)
}
