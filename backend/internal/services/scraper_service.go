package services

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

// ArticleMetadata represents the extracted metadata from a webpage
type ArticleMetadata struct {
	Title        string
	Content      string
	Author       string
	SiteName     string
	ThumbnailURL string
	PublishedAt  *time.Time
	Description  string
	Language     string
}

// ScraperService handles web scraping operations
type ScraperService struct {
	collector *colly.Collector
}

// NewScraperService creates a new scraper service
func NewScraperService() *ScraperService {
	c := colly.NewCollector(
		colly.Async(true),
		colly.AllowURLRevisit(),
	)

	// Set random user agent
	extensions.RandomUserAgent(c)

	// Set limits
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second,
	})

	// Set timeouts
	c.SetRequestTimeout(30 * time.Second)

	return &ScraperService{
		collector: c,
	}
}

// ExtractMetadata extracts metadata from the given URL
func (s *ScraperService) ExtractMetadata(targetURL string) (*ArticleMetadata, error) {
	// Validate URL
	_, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	metadata := &ArticleMetadata{
		Language: "ja", // Default to Japanese
	}

	// Set up callbacks for metadata extraction
	s.collector.OnHTML("head", func(e *colly.HTMLElement) {
		// Extract title
		if metadata.Title == "" {
			metadata.Title = e.ChildText("title")
			if metadata.Title == "" {
				metadata.Title = e.ChildAttr("meta[property='og:title']", "content")
			}
		}

		// Extract description
		metadata.Description = e.ChildAttr("meta[name='description']", "content")
		if metadata.Description == "" {
			metadata.Description = e.ChildAttr("meta[property='og:description']", "content")
		}

		// Extract thumbnail
		metadata.ThumbnailURL = e.ChildAttr("meta[property='og:image']", "content")
		if metadata.ThumbnailURL == "" {
			metadata.ThumbnailURL = e.ChildAttr("meta[name='twitter:image']", "content")
		}

		// Extract site name
		metadata.SiteName = e.ChildAttr("meta[property='og:site_name']", "content")

		// Extract author
		metadata.Author = e.ChildAttr("meta[name='author']", "content")
		if metadata.Author == "" {
			metadata.Author = e.ChildAttr("meta[property='article:author']", "content")
		}

		// Extract published date
		publishedTime := e.ChildAttr("meta[property='article:published_time']", "content")
		if publishedTime != "" {
			if t, err := time.Parse(time.RFC3339, publishedTime); err == nil {
				metadata.PublishedAt = &t
			}
		}

		// Extract language
		lang := e.ChildAttr("meta[property='og:locale']", "content")
		if lang != "" {
			metadata.Language = strings.Split(lang, "_")[0]
		}
	})

	// Extract main content
	s.collector.OnHTML("article", func(e *colly.HTMLElement) {
		content := e.Text
		if len(content) > len(metadata.Content) {
			metadata.Content = content
		}
	})

	// Fallback content extraction
	s.collector.OnHTML("main", func(e *colly.HTMLElement) {
		if metadata.Content == "" {
			metadata.Content = e.Text
		}
	})

	s.collector.OnHTML("body", func(e *colly.HTMLElement) {
		if metadata.Content == "" {
			// Try to find content in common content containers
			selectors := []string{
				".content",
				"#content",
				".post-content",
				".entry-content",
				".article-content",
				".main-content",
			}
			
			for _, selector := range selectors {
				content := e.ChildText(selector)
				if content != "" && len(content) > len(metadata.Content) {
					metadata.Content = content
					break
				}
			}
		}
	})

	// Error handling
	var scraperErr error
	s.collector.OnError(func(r *colly.Response, err error) {
		scraperErr = fmt.Errorf("scraping failed for %s: %w", targetURL, err)
	})

	// Visit the URL
	err = s.collector.Visit(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to visit URL: %w", err)
	}

	// Wait for async operations
	s.collector.Wait()

	if scraperErr != nil {
		return nil, scraperErr
	}

	// Clean up extracted content
	metadata.Content = cleanText(metadata.Content)
	metadata.Title = cleanText(metadata.Title)
	metadata.Description = cleanText(metadata.Description)

	return metadata, nil
}

// cleanText removes excessive whitespace and cleans up text
func cleanText(text string) string {
	// Remove extra whitespace
	text = strings.TrimSpace(text)
	
	// Replace multiple spaces with single space
	text = strings.Join(strings.Fields(text), " ")
	
	// Limit length for database storage
	if len(text) > 10000 {
		text = text[:10000] + "..."
	}
	
	return text
}

// ExtractContentForSite extracts content with site-specific rules
func (s *ScraperService) ExtractContentForSite(targetURL string) (*ArticleMetadata, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	// Site-specific extraction rules
	host := parsedURL.Host

	switch {
	case strings.Contains(host, "qiita.com"):
		return s.extractQiitaContent(targetURL)
	case strings.Contains(host, "zenn.dev"):
		return s.extractZennContent(targetURL)
	case strings.Contains(host, "note.com"):
		return s.extractNoteContent(targetURL)
	default:
		// Use generic extraction
		return s.ExtractMetadata(targetURL)
	}
}

// extractQiitaContent extracts content specifically from Qiita
func (s *ScraperService) extractQiitaContent(targetURL string) (*ArticleMetadata, error) {
	metadata := &ArticleMetadata{Language: "ja"}

	s.collector.OnHTML("article", func(e *colly.HTMLElement) {
		metadata.Title = e.ChildText("h1")
		metadata.Content = e.ChildText(".it-MdContent")
		metadata.Author = e.ChildText(".it-Header_authorName")
	})

	err := s.collector.Visit(targetURL)
	if err != nil {
		return nil, err
	}

	s.collector.Wait()
	return metadata, nil
}

// extractZennContent extracts content specifically from Zenn
func (s *ScraperService) extractZennContent(targetURL string) (*ArticleMetadata, error) {
	metadata := &ArticleMetadata{Language: "ja"}

	s.collector.OnHTML("article", func(e *colly.HTMLElement) {
		metadata.Title = e.ChildText("h1")
		metadata.Content = e.ChildText(".znc")
	})

	err := s.collector.Visit(targetURL)
	if err != nil {
		return nil, err
	}

	s.collector.Wait()
	return metadata, nil
}

// extractNoteContent extracts content specifically from Note
func (s *ScraperService) extractNoteContent(targetURL string) (*ArticleMetadata, error) {
	metadata := &ArticleMetadata{Language: "ja"}

	s.collector.OnHTML(".note-common-container", func(e *colly.HTMLElement) {
		metadata.Title = e.ChildText("h1")
		metadata.Content = e.ChildText(".note-body")
		metadata.Author = e.ChildText(".o-noteContentHeader__authorName")
	})

	err := s.collector.Visit(targetURL)
	if err != nil {
		return nil, err
	}

	s.collector.Wait()
	return metadata, nil
}