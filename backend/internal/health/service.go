package health

import (
	"context"
	"database/sql"
)

type Service struct {
	db *sql.DB
}

type Status struct {
	Status   string `json:"status"`
	Database string `json:"database,omitempty"`
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Check(ctx context.Context) Status {
	return Status{
		Status: "ok",
	}
}

func (s *Service) CheckDatabase(ctx context.Context) Status {
	if err := s.db.PingContext(ctx); err != nil {
		return Status{
			Status:   "error",
			Database: "disconnected",
		}
	}

	return Status{
		Status:   "ok",
		Database: "connected",
	}
}
