package model

import (
	"github.com/uptrace/bun"
	"payment-gateway/internal/app/domain"
	"time"
)

type Priority int

const (
	High Priority = iota
	Medium
	Low
)

func (p Priority) String() string {
	switch p {
	case High:
		return "high"
	case Medium:
		return "medium"
	case Low:
		return "low"
	default:
		return "unknown"
	}
}

func StrPriorityToInt(value string) (int, error) {
	switch value {
	case "high":
		return 0, nil
	case "medium":
		return 1, nil
	case "low":
		return 2, nil
	default:
		return 0, domain.ErrInvalidGatewayPriority
	}
}

func IntPriorityToString(value int) (string, error) {
	switch value {
	case 0:
		return "high", nil
	case 1:
		return "medium", nil
	case 2:
		return "low", nil
	default:
		return "", domain.ErrInvalidGatewayPriority
	}
}

type Gateway struct {
	bun.BaseModel       `bun:"table:gateways"`
	ID                  int       `bun:",pk,autoincrement"`
	Name                string    `bun:",unique,notnull"`
	DataFormatSupported string    `bun:",notnull"`
	Priority            int       `bun:",notnull"`
	CreatedAt           time.Time `bun:",default:current_timestamp"`
	UpdatedAt           time.Time `bun:",default:current_timestamp"`
}
