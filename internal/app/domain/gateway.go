package domain

import "time"

type Gateway struct {
	id                  int
	name                string
	dataFormatSupported string
	priority            string
	createdAt           time.Time
	updatedAt           time.Time
}

type GatewayData struct {
	ID                  int
	Name                string
	DataFormatSupported string
	Priority            string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func NewGateway(data GatewayData) (Gateway, error) {
	return Gateway{
		id:                  data.ID,
		name:                data.Name,
		priority:            data.Priority,
		dataFormatSupported: data.DataFormatSupported,
		createdAt:           data.CreatedAt,
		updatedAt:           data.UpdatedAt,
	}, nil
}

func (g Gateway) ID() int {
	return g.id
}

func (g Gateway) Name() string {
	return g.name
}

func (g Gateway) DataFormatSupported() string {
	return g.dataFormatSupported
}

func (g Gateway) CreatedAt() time.Time {
	return g.createdAt
}

func (g Gateway) UpdatedAt() time.Time {
	return g.updatedAt
}
