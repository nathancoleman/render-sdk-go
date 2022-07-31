package render

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return e.ID + ": " + e.Message
}

func ErrorFromResponse(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}

	e := &ErrorResponse{}
	if err := json.NewDecoder(resp.Body).Decode(e); err != nil {
		return nil
	}
	return e
}
