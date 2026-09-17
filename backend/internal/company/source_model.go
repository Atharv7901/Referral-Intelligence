package company

import "time"

type SourceType string

const (
	SourceGreenhouse SourceType = "greenhouse"
	SourceLever      SourceType = "lever"
	SourceAshby      SourceType = "ashby"
	SourceCareerPage SourceType = "career_page"
)

type CompanySource struct {
	ID         uint64     `json:"id"`
	CompanyID  uint64     `json:"company_id"`
	Source     SourceType `json:"source"`
	ExternalID *string    `json:"external_id,omitempty"`
	SourceURL  *string    `json:"source_url,omitempty"`
	Metadata   []byte     `json:"metadata,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
