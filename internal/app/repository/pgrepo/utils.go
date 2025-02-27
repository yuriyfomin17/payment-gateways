package pgrepo

import (
	"payment-gateway/internal/app/domain"
	"payment-gateway/internal/app/repository/model"
)

func domainToTransaction(transaction domain.Transaction) (model.Transaction, error) {
	return model.Transaction{
		ID:        transaction.ID(),
		Amount:    transaction.Amount(),
		Type:      transaction.Type(),
		Status:    transaction.Status(),
		CreatedAt: transaction.CreatedAt(),
		GatewayID: transaction.GatewayID(),
		CountryID: transaction.CountryID(),
		UserID:    transaction.UserID(),
	}, nil
}

func gatewayToDomain(gateway model.Gateway) (domain.Gateway, error) {
	strPriority, err := model.IntPriorityToString(gateway.Priority)
	if err != nil {
		return domain.Gateway{}, domain.ErrInvalidGatewayPriority
	}
	return domain.NewGateway(domain.GatewayData{
		ID:                  gateway.ID,
		Name:                gateway.Name,
		Priority:            strPriority,
		DataFormatSupported: gateway.DataFormatSupported,
		CreatedAt:           gateway.CreatedAt,
		UpdatedAt:           gateway.UpdatedAt,
	})
}

func userToDomain(user model.User) (domain.User, error) {
	return domain.NewUser(domain.UserData{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CountryID: user.CountryID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func countryToDomain(country model.Country) (domain.Country, error) {
	return domain.NewCountry(domain.CountryData{
		ID:        country.ID,
		Name:      country.Name,
		Code:      country.Code,
		Currency:  country.Currency,
		CreatedAt: country.CreatedAt,
		UpdatedAt: country.UpdatedAt,
	})
}
