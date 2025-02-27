package pgrepo

import (
	"context"
	"database/sql"
	"errors"
	"payment-gateway/internal/app/domain"
	"payment-gateway/internal/app/repository/model"
	"payment-gateway/internal/pkg/pg"
)

type UserRepo struct {
	db *pg.DB
}

func NewUserRepo(db *pg.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (repo UserRepo) GetUserByID(ctx context.Context, userID int) (domain.User, error) {
	var user model.User
	err := repo.db.NewSelect().Model(&user).
		Where("id = ?", userID).
		Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrNotFound
	}
	return userToDomain(user)
}
