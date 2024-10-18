package openai

import (
	"net/http"

	"github.com/thalesfsp/customerror"
)

// ProcessResponse processes the response from the API.
func ProcessResponse(resp ResponseBody) (string, error) {
	for _, choice := range resp.Choices {
		if choice.Message.Content != "" {
			return choice.Message.Content, nil
		}
	}

	return "", customerror.New(
		"no content",
		customerror.WithStatusCode(http.StatusNoContent),
	)
}
