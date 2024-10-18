package ollama

import (
	"net/http"

	"github.com/thalesfsp/customerror"
)

// ProcessResponse processes the response from the API.
func ProcessResponse(response ResponseBody) (string, error) {
	if response.Message.Content == "" {
		return "", customerror.New(
			"no content",
			customerror.WithStatusCode(http.StatusNoContent),
		)
	}

	return response.Message.Content, nil
}
