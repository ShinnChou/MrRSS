package feed

import (
	"MrRSS/internal/models"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

// processArticles processes RSS feed items and converts them to Article models
func (f *Fetcher) processArticles(feed models.Feed, items []*gofeed.Item) []*models.Article {
	// Check translation settings
	translationEnabledStr, _ := f.db.GetSetting("translation_enabled")
	targetLang, _ := f.db.GetSetting("target_language")
	translationEnabled := translationEnabledStr == "true"

	var articles []*models.Article

	for _, item := range items {
		published := time.Now()
		if item.PublishedParsed != nil {
			published = *item.PublishedParsed
		}

		imageURL := extractImageURL(item)
		audioURL := extractAudioURL(item)

		translatedTitle := ""
		if translationEnabled && f.translator != nil {
			t, err := f.translator.Translate(item.Title, targetLang)
			if err == nil {
				translatedTitle = t
			}
		}

		// Extract content from RSS item (prefer Content over Description)
		content := item.Content
		if content == "" {
			content = item.Description
		}

		// Generate title if missing
		title := item.Title
		if title == "" {
			title = generateTitleFromContent(content)
		}

		article := &models.Article{
			FeedID:          feed.ID,
			Title:           title,
			URL:             item.Link,
			ImageURL:        imageURL,
			AudioURL:        audioURL,
			Content:         content,
			PublishedAt:     published,
			TranslatedTitle: translatedTitle,
		}
		articles = append(articles, article)
	}

	return articles
}

// extractImageURL extracts the image URL from a feed item
func extractImageURL(item *gofeed.Item) string {
	// Try item.Image first
	if item.Image != nil {
		return item.Image.URL
	}

	// Try enclosures for images (check various image MIME types)
	for _, enc := range item.Enclosures {
		if strings.HasPrefix(enc.Type, "image/") {
			return enc.URL
		}
	}

	// Fallback: Try to find image in description/content
	content := item.Content
	if content == "" {
		content = item.Description
	}

	re := regexp.MustCompile(`<img[^>]+src="([^">]+)"`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// extractAudioURL extracts the audio URL from a feed item (for podcasts)
func extractAudioURL(item *gofeed.Item) string {
	// Try enclosures for audio files
	for _, enc := range item.Enclosures {
		// Check for audio MIME types
		if strings.HasPrefix(enc.Type, "audio/") {
			return enc.URL
		}
	}

	return ""
}

// generateTitleFromContent generates a title from content when title is missing
func generateTitleFromContent(content string) string {
	if content == "" {
		return "Untitled Article"
	}

	// Remove HTML tags
	htmlTagRegex := regexp.MustCompile(`<[^>]+>`)
	plainText := htmlTagRegex.ReplaceAllString(content, "")

	// Trim whitespace
	plainText = strings.TrimSpace(plainText)

	// Limit to 100 characters
	if len(plainText) > 100 {
		plainText = plainText[:100] + "..."
	}

	// If still empty after cleaning, use default
	if plainText == "" {
		return "Untitled Article"
	}

	return plainText
}
