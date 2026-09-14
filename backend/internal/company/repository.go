package company

import (
	"context"
	"database/sql"
)

type Repository interface {
	Create(ctx context.Context, company *Company) error
	GetById(ctx context.Context, id uint64) (*Company, error)
	GetBySlug(ctx context.Context, slug string) (*Company, error)
	List(ctx context.Context, search string) ([]Company, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, company *Company) error {
	query := `
		INSERT INTO companies (name, slug, website, linkedin_url, career_page_url, company_type, is_active)
		VALUES (?,?,?,?,?,?,?)
	`

	result, err := r.db.ExecContext(ctx, query, company.Name, company.Slug, company.Website, company.LinkedInUrl, company.CareerPageUrl, company.CompanyType, company.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	company.Id = uint64(id)

	return nil
}

func (r *repository) GetById(ctx context.Context, id uint64) (*Company, error) {
	query := `
		SELECT id, name, slug, website, linkedin_url, career_page_url, company_type, is_active, created_at, updated_at
		FROM companies WHERE id = ?
	`

	var company Company

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&company.Id,
		&company.Name,
		&company.Slug,
		&company.Website,
		&company.LinkedInUrl,
		&company.CareerPageUrl,
		&company.CompanyType,
		&company.IsActive,
		&company.CreatedAt,
		&company.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (r *repository) GetBySlug(ctx context.Context, slug string) (*Company, error) {
	query := `
		SELECT
			id,
			name,
			slug,
			website,
			linkedin_url,
			career_page_url,
			company_type,
			is_active,
			created_at,
			updated_at
		FROM companies
		WHERE slug = ?
	`

	var company Company

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&company.Id,
		&company.Name,
		&company.Slug,
		&company.Website,
		&company.LinkedInUrl,
		&company.CareerPageUrl,
		&company.CompanyType,
		&company.IsActive,
		&company.CreatedAt,
		&company.UpdateAt,
	)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (r *repository) List(ctx context.Context, search string) ([]Company, error) {
	query := `
		SELECT id, name, slug, website, linkedin_url, career_page_url, company_type, is_active, created_at, updated_at
		FROM companies WHERE is_active = TRUE
	`

	args := []interface{}{}

	if search != "" {
		query += `
			AND name LIKE ?
		`
		args = append(args, "%"+search+"%")
	}

	query += `
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companies := make([]Company, 0)

	for rows.Next() {
		var company Company

		err := rows.Scan(
			&company.Id,
			&company.Name,
			&company.Slug,
			&company.Website,
			&company.LinkedInUrl,
			&company.CareerPageUrl,
			&company.CompanyType,
			&company.IsActive,
			&company.CreatedAt,
			&company.UpdateAt,
		)
		if err != nil {
			return nil, err
		}

		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}
