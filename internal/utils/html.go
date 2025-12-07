package utils

import (
	"regexp"
	"strings"
)

// CleanHTML sanitizes HTML content by fixing common malformed patterns
// that can cause rendering issues.
func CleanHTML(html string) string {
	if html == "" {
		return html
	}

	// Fix malformed opening tags like <p--> to <p>
	// This pattern matches tags like <p-->, <div-->, etc.
	malformedTagRegex := regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)\s*--+>`)
	html = malformedTagRegex.ReplaceAllString(html, "<$1>")

	// Fix malformed self-closing tags like <img-->, <br--> to <img>, <br>
	// Some feeds have broken self-closing tags with or without attributes
	// Pattern 1: Tags with attributes (e.g., <img src="..." -->)
	malformedSelfClosingWithAttrs := regexp.MustCompile(`<(img|br|hr|input|meta|link)\s+([^>]+?)--+>`)
	html = malformedSelfClosingWithAttrs.ReplaceAllString(html, "<$1 $2>")
	
	// Pattern 2: Tags without attributes (e.g., <br-->)
	malformedSelfClosingNoAttrs := regexp.MustCompile(`<(img|br|hr|input|meta|link)\s*--+>`)
	html = malformedSelfClosingNoAttrs.ReplaceAllString(html, "<$1>")

	// Trim any leading/trailing whitespace
	html = strings.TrimSpace(html)

	return html
}
