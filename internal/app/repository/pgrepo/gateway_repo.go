package pgrepo

import (
	"context"
	"fmt"
	"payment-gateway/internal/app/domain"
	"payment-gateway/internal/app/repository/model"
	"payment-gateway/internal/pkg/pg"
	"time"
)

type GatewayRepo struct {
	db *pg.DB
}

func NewGatewayRepo(db *pg.DB) *GatewayRepo {
	return &GatewayRepo{
		db: db,
	}
}

func (repo GatewayRepo) UpdateGatewayPriority(ctx context.Context, gatewayID int64, priority string) error {
	intPriority, err := model.StrPriorityToInt(priority)
	if err != nil {
		return fmt.Errorf("failed to update priority for gateway with id %d: %w", gatewayID, err)
	}
	_, err = repo.db.NewUpdate().
		Model(&model.Gateway{}).
		Set("priority = ?", intPriority).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", gatewayID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update priority for gateway with id %d: %w", gatewayID, err)
	}
	return nil
}

func (repo GatewayRepo) GetSupportedGatewaysByCountrySortedByPriorities(ctx context.Context, countryId int, dataFormat string) ([]domain.Gateway, error) {
	var modelGateways []model.Gateway
	err := repo.db.NewSelect().
		Model(&modelGateways).
		Where(`gateway.id IN (SELECT gateway_id FROM gateway_countries WHERE country_id = ? ) AND gateway.data_format_supported LIKE '%' || ? || '%'`, countryId, dataFormat).
		Order("priority ASC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve countries for gateway %d: %w", countryId, err)
	}

	domainCountries := make([]domain.Gateway, len(modelGateways))
	for i, gateway := range modelGateways {
		domainCountry, err := gatewayToDomain(gateway)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain gateway: %w", err)
		}
		domainCountries[i] = domainCountry
	}
	return domainCountries, nil
}
