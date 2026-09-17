package company

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type SourceResolver interface {
	Resolve(ctx context.Context, company *Company) ([]CompanySource, error)
}

type sourceResolver struct {
	repository SourceRepository
	detector   ATSDetector
}

func NewSourceResolver(
	repository SourceRepository,
	detector ATSDetector,
) sourceResolver {
	return sourceResolver{
		repository: repository,
		detector:   detector,
	}
}

func (r *sourceResolver) Resolve(
	ctx context.Context,
	company *Company,
) ([]CompanySource, error) {
	if company == nil {
		return nil, fmt.Errorf("company cannot be nil")
	}

	if company.CareerPageUrl == nil ||
		strings.TrimSpace(*company.CareerPageUrl) == "" {
		return nil, nil
	}
	fmt.Println("reached here")
	detected, err := r.detector.Detect(
		ctx,
		strings.TrimSpace(*company.CareerPageUrl),
	)
	if err != nil {
		return nil, fmt.Errorf("detect ATS: %w", err)
	}

	if detected == nil {
		source := CompanySource{
			CompanyID: company.Id,
			Source:    SourceCareerPage,
			SourceURL: company.CareerPageUrl,
		}

		if err := r.repository.Upsert(ctx, &source); err != nil {
			return nil, fmt.Errorf(
				"persist company source: %w",
				err,
			)
		}

		return []CompanySource{source}, nil
	}

	metadata, err := json.Marshal(map[string]any{
		"confidence":    detected.Confidence,
		"detected_from": "career_page",
	})
	if err != nil {
		return nil, fmt.Errorf("marshal source metadata: %w", err)
	}

	source := CompanySource{
		CompanyID: company.Id,
		Source:    detected.Type,
		ExternalID: stringPtr(
			detected.ExternalID,
		),
		SourceURL: stringPtr(
			detected.SourceURL,
		),
		Metadata: metadata,
	}

	fmt.Println("this is the source", source)

	if err := r.repository.Upsert(ctx, &source); err != nil {
		return nil, fmt.Errorf(
			"persist company source: %w",
			err,
		)
	}
	return []CompanySource{source}, nil
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
