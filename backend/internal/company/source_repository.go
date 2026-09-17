package company

import (
	"context"
	"database/sql"
	"fmt"
)

type SourceRepository interface {
	Create(ctx context.Context, source *CompanySource) error
	GetByCompanyID(ctx context.Context, companyID uint64) ([]CompanySource, error)
	Upsert(ctx context.Context, source *CompanySource) error
	GetByCompanyAndType(
		ctx context.Context,
		companyID uint64,
		sourceType SourceType,
	) (*CompanySource, error)
}

type sourceRepository struct {
	db *sql.DB
}

func NewSourceRepository(db *sql.DB) sourceRepository {
	return sourceRepository{
		db: db,
	}
}

func (r *sourceRepository) Create(
	ctx context.Context,
	source *CompanySource,
) error {
	query := `
		INSERT INTO company_sources (
			company_id,
			source,
			external_id,
			source_url,
			metadata
		)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		source.CompanyID,
		source.Source,
		source.ExternalID,
		source.SourceURL,
		source.Metadata,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	source.ID = uint64(id)

	return nil
}

func (r *sourceRepository) GetByCompanyID(
	ctx context.Context,
	companyID uint64,
) ([]CompanySource, error) {
	query := `
		SELECT
			id,
			company_id,
			source,
			external_id,
			source_url,
			metadata,
			created_at,
			updated_at
		FROM company_sources
		WHERE company_id = ?
		ORDER BY id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sources := make([]CompanySource, 0)

	for rows.Next() {
		var source CompanySource

		err := rows.Scan(
			&source.ID,
			&source.CompanyID,
			&source.Source,
			&source.ExternalID,
			&source.SourceURL,
			&source.Metadata,
			&source.CreatedAt,
			&source.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}

func (r *sourceRepository) GetByCompanyAndType(
	ctx context.Context,
	companyID uint64,
	sourceType SourceType,
) (*CompanySource, error) {
	query := `
		SELECT
			id,
			company_id,
			source,
			external_id,
			source_url,
			metadata,
			created_at,
			updated_at
		FROM company_sources
		WHERE company_id = ?
		  AND source = ?
		LIMIT 1
	`

	var source CompanySource

	err := r.db.QueryRowContext(
		ctx,
		query,
		companyID,
		sourceType,
	).Scan(
		&source.ID,
		&source.CompanyID,
		&source.Source,
		&source.ExternalID,
		&source.SourceURL,
		&source.Metadata,
		&source.CreatedAt,
		&source.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &source, nil
}

func (r *sourceRepository) Upsert(
	ctx context.Context,
	source *CompanySource,
) error {
	const query = `
		INSERT INTO company_sources (
			company_id,
			source,
			external_id,
			source_url,
			metadata
		)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			company_id = VALUES(company_id),
			source_url = VALUES(source_url),
			metadata = VALUES(metadata),
			updated_at = CURRENT_TIMESTAMP
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		source.CompanyID,
		source.Source,
		source.ExternalID,
		source.SourceURL,
		source.Metadata,
	)

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	source.ID = uint64(id)

	if err != nil {
		return fmt.Errorf(
			"upsert company source: %w",
			err,
		)
	}

	return nil
}
