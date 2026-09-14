package company

import (
	"context"
	"errors"
	"strings"
)

type Service interface {
	Create(ctx context.Context, company *Company) error
	GetByID(ctx context.Context, id uint64) (*Company, error)
	GetBySlug(ctx context.Context, slug string) (*Company, error)
	List(ctx context.Context, search string) ([]Company, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) service {
	return service{
		repository: repository,
	}
}

var (
	ErrInvalidCompanyName = errors.New("company name is required")
	ErrInvalidCompanySlug = errors.New("company slug is required")
	ErrInvalidCompanyType = errors.New("company type is required")
)

func (s *service) Create(ctx context.Context, company *Company) error {
	if company == nil {
		return errors.New("company is required")
	}

	company.Name = strings.TrimSpace(company.Name)
	company.Slug = strings.TrimSpace(strings.ToLower(company.Slug))
	company.CompanyType = strings.TrimSpace(company.CompanyType)

	if company.Name == "" {
		return ErrInvalidCompanyName
	}

	if company.Slug == "" {
		return ErrInvalidCompanySlug
	}

	if company.CompanyType == "" {
		return ErrInvalidCompanyType
	}

	return s.repository.Create(ctx, company)
}

func (s *service) GetByID(ctx context.Context, id uint64) (*Company, error) {
	if id == 0 {
		return nil, errors.New("invalid company id")
	}

	return s.repository.GetById(ctx, id)
}

func (s *service) GetBySlug(ctx context.Context, slug string) (*Company, error) {
	slug = strings.TrimSpace(strings.ToLower(slug))

	if slug == "" {
		return nil, ErrInvalidCompanySlug
	}

	return s.repository.GetBySlug(ctx, slug)
}

func (s *service) List(ctx context.Context, search string) ([]Company, error) {
	search = strings.TrimSpace(search)

	return s.repository.List(ctx, search)
}
