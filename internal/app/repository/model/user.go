package model

import (
	"github.com/uptrace/bun"
	_ "github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            int       `bun:",pk,autoincrement"`
	Username      string    `bun:",unique,notnull"`
	Email         string    `bun:",unique,notnull"`
	Password      string    `bun:",notnull"`
	CountryID     int       `bun:",nullzero"`
	CreatedAt     time.Time `bun:",default:current_timestamp"`
	UpdatedAt     time.Time `bun:",default:current_timestamp"`
}
