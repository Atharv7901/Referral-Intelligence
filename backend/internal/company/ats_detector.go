package company

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type httpATSDetector struct {
	client *http.Client
}

func NewATSDetector() httpATSDetector {
	return httpATSDetector{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (d *httpATSDetector) Detect(
	ctx context.Context,
	careerPageURL string,
) (*DetectedSource, error) {
	parsedURL, err := url.Parse(careerPageURL)
	if err != nil {
		return nil, fmt.Errorf("parse career page URL: %w", err)
	}

	// First try the URL itself.
	if source := detectFromURL(parsedURL); source != nil {
		return source, nil
	}
	fmt.Println("going here 1")
	// For custom career pages, inspect the page and follow redirects.
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		careerPageURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create ATS detection request: %w", err)
	}

	req.Header.Set(
		"User-Agent",
		"ReferralIntelligenceBot/1.0",
	)
	req.Header.Set(
		"Accept",
		"text/html,application/xhtml+xml",
	)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch career page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, nil
	}

	// The HTTP client follows redirects automatically.
	if source := detectFromURL(resp.Request.URL); source != nil {
		return source, nil
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("read career page: %w", err)
	}

	return detectFromHTML(string(body))
}

func detectFromURL(u *url.URL) *DetectedSource {
	if u == nil {
		return nil
	}

	host := strings.ToLower(u.Hostname())
	path := strings.Trim(u.Path, "/")

	switch {
	case isGreenhouseHost(host):
		return &DetectedSource{
			Type:       SourceGreenhouse,
			ExternalID: firstPathSegment(path),
			SourceURL:  u.String(),
			Confidence: 1.0,
		}

	case isLeverHost(host):
		return &DetectedSource{
			Type:       SourceLever,
			ExternalID: firstPathSegment(path),
			SourceURL:  u.String(),
			Confidence: 1.0,
		}

	case isAshbyHost(host):
		return &DetectedSource{
			Type:       SourceAshby,
			ExternalID: firstPathSegment(path),
			SourceURL:  u.String(),
			Confidence: 1.0,
		}
	}

	return nil
}

func detectFromHTML(html string) (*DetectedSource, error) {
	lowerHTML := strings.ToLower(html)

	// Prefer explicit ATS URLs when available.

	if strings.Contains(lowerHTML, "boards.greenhouse.io/") ||
		strings.Contains(lowerHTML, "job-boards.greenhouse.io/") {
		if id := extractGreenhouseID(lowerHTML); id != "" {
			return &DetectedSource{
				Type:       SourceGreenhouse,
				ExternalID: id,
				SourceURL:  "https://boards.greenhouse.io/" + id,
				Confidence: 0.9,
			}, nil
		}
	}

	if strings.Contains(lowerHTML, "jobs.lever.co/") ||
		strings.Contains(lowerHTML, "jobs.eu.lever.co/") {
		if id := extractLeverID(lowerHTML); id != "" {
			return &DetectedSource{
				Type:       SourceLever,
				ExternalID: id,
				SourceURL:  "https://jobs.lever.co/" + id,
				Confidence: 0.9,
			}, nil
		}
	}

	if strings.Contains(lowerHTML, "jobs.ashbyhq.com/") {
		if id := extractAshbyID(lowerHTML); id != "" {
			return &DetectedSource{
				Type:       SourceAshby,
				ExternalID: id,
				SourceURL:  "https://jobs.ashbyhq.com/" + id,
				Confidence: 0.9,
			}, nil
		}
	}

	// We found evidence of an ATS but could not reliably
	// extract the board identifier. Do not create a bad source.
	return nil, nil
}

func isGreenhouseHost(host string) bool {
	return host == "boards.greenhouse.io" ||
		host == "job-boards.greenhouse.io"
}

func isLeverHost(host string) bool {
	return host == "jobs.lever.co" ||
		host == "jobs.eu.lever.co"
}

func isAshbyHost(host string) bool {
	return host == "jobs.ashbyhq.com"
}

func firstPathSegment(path string) string {
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")

	for _, part := range parts {
		if part != "" {
			return part
		}
	}

	return ""
}

func extractGreenhouseID(html string) string {
	return extractIDAfter(
		html,
		"boards.greenhouse.io/",
		"job-boards.greenhouse.io/",
	)
}

func extractLeverID(html string) string {
	return extractIDAfter(
		html,
		"jobs.lever.co/",
		"jobs.eu.lever.co/",
	)
}

func extractAshbyID(html string) string {
	return extractIDAfter(
		html,
		"jobs.ashbyhq.com/",
	)
}

func extractIDAfter(html string, prefixes ...string) string {
	for _, prefix := range prefixes {
		index := strings.Index(html, prefix)

		if index == -1 {
			continue
		}

		value := html[index+len(prefix):]

		end := len(value)

		for i, char := range value {
			if char == '/' ||
				char == '"' ||
				char == '\'' ||
				char == '?' ||
				char == '#' ||
				char == ' ' {
				end = i
				break
			}
		}

		id := strings.TrimSpace(value[:end])

		if id != "" {
			return id
		}
	}

	return ""
}
