package anthropic

import (
	"net/http"

	"github.com/thalesfsp/customerror"
)

// ProcessResponse processes the response from the API.
func ProcessResponse(resp ResponseBody) (string, error) {
	for _, content := range resp.Content {
		if content.Text != "" {
			return content.Text, nil
		}
	}

	return "", customerror.New(
		"no content",
		customerror.WithStatusCode(http.StatusNoContent),
	)
}
