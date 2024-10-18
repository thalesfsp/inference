package message

// Message to be sent to the LLM provider.
type Message struct {
	Content string `json:"content"`
	Role    Role   `json:"role"`
}

// NewMessages creates a new list of messages properly separated by role where
// system messages come first and user messages come last.
func NewMessages(
	systemMessages []string,
	userMessages []string,
) []Message {
	// Stores the final messages.
	finalMessages := []Message{}

	for _, systemMessage := range systemMessages {
		finalMessages = append(finalMessages, Message{
			Content: systemMessage,
			Role:    System,
		})
	}

	for _, userMessage := range userMessages {
		finalMessages = append(finalMessages, Message{
			Content: userMessage,
			Role:    User,
		})
	}

	return finalMessages
}
