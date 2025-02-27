package domain

import "time"

type Country struct {
	id        int
	name      string
	code      string
	currency  string
	createdAt time.Time
	updatedAt time.Time
}

type CountryData struct {
	ID        int
	Name      string
	Code      string
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCountry(data CountryData) (Country, error) {
	return Country{
		id:        data.ID,
		name:      data.Name,
		code:      data.Code,
		currency:  data.Currency,
		createdAt: data.CreatedAt,
		updatedAt: data.UpdatedAt,
	}, nil
}

func (c Country) ID() int {
	return c.id
}

func (c Country) Name() string {
	return c.name
}

func (c Country) Code() string {
	return c.code
}

func (c Country) Currency() string {
	return c.currency
}

func (c Country) CreatedAt() time.Time {
	return c.createdAt
}

func (c Country) UpdatedAt() time.Time {
	return c.updatedAt
}
