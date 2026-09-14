package company

import "time"

type Company struct {
	Id            uint64    `json:"id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Website       *string   `json:"website,omitempty"`
	LinkedInUrl   *string   `json:"linkedin_url,omitempty"`
	CareerPageUrl *string   `json:"career_page_url,omitempty"`
	CompanyType   string    `json:"company_type"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdateAt      time.Time `json:"updated_at"`
}
