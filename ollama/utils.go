package ollama

import (
	"strings"

	"github.com/thalesfsp/customerror"
)

// ProcessResponse processes the response from the API.
func ProcessResponse(response ResponseBody) (string, error) {
	if len(strings.TrimSpace(response.Message.Content)) == 0 {
		return "", customerror.NewMissingError("content")
	}

	return response.Message.Content, nil
}
