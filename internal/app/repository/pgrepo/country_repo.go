package pgrepo

import (
	"context"
	"database/sql"
	"errors"
	"payment-gateway/internal/app/domain"
	"payment-gateway/internal/app/repository/model"
	"payment-gateway/internal/pkg/pg"
)

type CountryRepo struct {
	db *pg.DB
}

func NewCountryRepo(db *pg.DB) *CountryRepo {
	return &CountryRepo{
		db: db,
	}
}

func (repo CountryRepo) GetCountryByID(ctx context.Context, countryID int, currencyId string) (domain.Country, error) {
	var country model.Country
	err := repo.db.NewSelect().Model(&country).
		Where("id = ? AND currency = ?", countryID, currencyId).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Country{}, domain.ErrNotFound
	}
	return countryToDomain(country)
}
