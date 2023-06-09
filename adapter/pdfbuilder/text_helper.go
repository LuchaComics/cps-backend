package pdfbuilder

import "strings"

func splitText(text string, maxChars int) []string {
	var lines []string
	words := strings.Fields(text)

	currentLine := ""
	for _, word := range words {
		if len(currentLine)+len(word) <= maxChars {
			// Word fits within the current line
			if currentLine != "" {
				currentLine += " " + word
			} else {
				currentLine += word
			}
		} else {
			// Word needs to be moved to the next line
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	// Add the last line
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

func getElementAtIndex(slice []string, index int) (string, bool) {
	if index < 0 || index >= len(slice) {
		return "", false
	}
	return slice[index], true
}
