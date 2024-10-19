package openai

import (
	"net/http"
	"strings"

	"github.com/thalesfsp/customerror"
)

// ProcessResponse processes the response from the API.
func ProcessResponse(resp ResponseBody) (string, error) {
	for _, choice := range resp.Choices {
		if len(strings.TrimSpace(choice.Message.Content)) == 0 {
			return choice.Message.Content, nil
		}
	}

	return "", customerror.New(
		"no content",
		customerror.WithStatusCode(http.StatusNoContent),
	)
}
