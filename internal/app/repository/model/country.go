package model

import (
	"github.com/uptrace/bun"
	"time"
)

type Country struct {
	bun.BaseModel `bun:"table:countries"`
	ID            int       `bun:",pk,autoincrement"`
	Name          string    `bun:",unique,notnull"`
	Code          string    `bun:",unique,notnull"`
	Currency      string    `bun:",notnull"`
	CreatedAt     time.Time `bun:",default:current_timestamp"`
	UpdatedAt     time.Time `bun:",default:current_timestamp"`
}
