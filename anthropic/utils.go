package anthropic

import (
	"net/http"
	"strings"

	"github.com/thalesfsp/customerror"
)

// ProcessResponse processes the response from the API.
func ProcessResponse(resp ResponseBody) (string, error) {
	for _, content := range resp.Content {
		if len(strings.TrimSpace(content.Text)) == 0 {
			return content.Text, nil
		}
	}

	return "", customerror.New(
		"no content",
		customerror.WithStatusCode(http.StatusNoContent),
	)
}
