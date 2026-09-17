package company

import "context"

type DetectedSource struct {
	Type       SourceType
	ExternalID string
	SourceURL  string
	Confidence float64
}

type ATSDetector interface {
	Detect(ctx context.Context, careerPageURL string) (*DetectedSource, error)
}
