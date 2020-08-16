package app

// TrimText truncates a piece of text by a specific length.
func TrimText(text string, truncLength int) string {
	if len(text) > truncLength {
		// Split string by rune length.
		// Ref: https://stackoverflow.com/a/46416046
		// TODO: Not the best solution. Consider Adrian's approach in the SO answer.
		return string([]rune(text)[:truncLength-3]) + "..."
	}
	return text
}
