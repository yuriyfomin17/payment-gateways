package services

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"payment-gateway/internal/app/domain"
	"payment-gateway/internal/app/services/mocks"
	"testing"
)

func TestFetchTxId_Success(t *testing.T) {
	// Given
	mockTransactionRepo := mocks.NewTransactionRepository(t)
	mockTransactionRepo.On("FetchTxId", mock.Anything).Return(int64(12345), nil)

	service := NewUserService(nil, mockTransactionRepo, nil, nil, nil, nil)
	ctx := context.Background()

	// When
	txId, err := service.FetchTxId(ctx)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, int64(12345), txId)
	mockTransactionRepo.AssertExpectations(t)
}

func TestFetchTxId_Error(t *testing.T) {
	// Given
	mockTransactionRepo := mocks.NewTransactionRepository(t)
	mockTransactionRepo.On("FetchTxId", mock.Anything).Return(int64(0), errors.New("some error"))

	service := NewUserService(nil, mockTransactionRepo, nil, nil, nil, nil)
	ctx := context.Background()

	// When
	txId, err := service.FetchTxId(ctx)

	// Then
	assert.Error(t, err)
	assert.Equal(t, int64(0), txId)
	assert.Equal(t, domain.ErrTransactionNotCreated, err)
	mockTransactionRepo.AssertExpectations(t)
}

func TestExecuteTransaction_Success(t *testing.T) {
	// Given
	mockTransactionRepo := mocks.NewTransactionRepository(t)
	mockUserRepo := mocks.NewUserRepository(t)
	mockCountryRepo := mocks.NewCountryRepository(t)
	mockGatewayRepo := mocks.NewGatewayRepository(t)
	mockFaultTolerance := mocks.NewFaultTolerance(t)

	ctx := context.Background()

	txData := domain.TransactionData{
		UserID:     1,
		Currency:   "USD",
		DataFormat: "json",
	}

	mockUser := domain.User{}
	mockUserRepo.On("GetUserByID", ctx, 1).Return(mockUser, nil)

	mockCountry := domain.CountryData{ID: 12}
	country, _ := domain.NewCountry(mockCountry)
	mockCountryRepo.On("GetCountryByID", ctx, mockUser.CountryID(), "USD").Return(country, nil)

	mockGateway, _ := domain.NewGateway(domain.GatewayData{ID: 100})
	mockGatewayRepo.On("GetSupportedGatewaysByCountrySortedByPriorities", ctx, 12, "json").Return([]domain.Gateway{mockGateway}, nil)

	mockFaultTolerance.On("RetryOperation", mock.Anything, 3).Return(nil)

	service := NewUserService(mockUserRepo, mockTransactionRepo, mockCountryRepo, mockGatewayRepo, mockFaultTolerance, nil)

	// When
	err := service.ExecuteTransaction(ctx, txData)

	// Then
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockCountryRepo.AssertExpectations(t)
	mockGatewayRepo.AssertExpectations(t)
	mockFaultTolerance.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestExecuteTransaction_UserNotFound(t *testing.T) {
	// Given
	mockUserRepo := mocks.NewUserRepository(t)
	mockTransactionRepo := mocks.NewTransactionRepository(t)

	ctx := context.Background()
	txData := domain.TransactionData{UserID: 1}

	mockUserRepo.On("GetUserByID", ctx, 1).Return(domain.User{}, errors.New("not found"))

	service := NewUserService(mockUserRepo, mockTransactionRepo, nil, nil, nil, nil)

	// When
	err := service.ExecuteTransaction(ctx, txData)

	// Then
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	mockUserRepo.AssertExpectations(t)
}

func TestExecuteTransaction_TransactionNotCreated(t *testing.T) {
	// Given
	mockTransactionRepo := mocks.NewTransactionRepository(t)
	mockUserRepo := mocks.NewUserRepository(t)
	mockCountryRepo := mocks.NewCountryRepository(t)
	mockGatewayRepo := mocks.NewGatewayRepository(t)
	mockFaultTolerance := mocks.NewFaultTolerance(t)

	ctx := context.Background()

	txData := domain.TransactionData{
		UserID:     1,
		Currency:   "USD",
		DataFormat: "json",
	}

	mockUser := domain.User{}
	mockUserRepo.On("GetUserByID", ctx, 1).Return(mockUser, nil)

	mockCountry := domain.CountryData{ID: 12}
	country, _ := domain.NewCountry(mockCountry)
	mockCountryRepo.On("GetCountryByID", ctx, mockUser.CountryID(), "USD").Return(country, nil)

	mockGateway, _ := domain.NewGateway(domain.GatewayData{ID: 100})
	mockGatewayRepo.On("GetSupportedGatewaysByCountrySortedByPriorities", ctx, 12, "json").Return([]domain.Gateway{mockGateway}, nil)

	// Simulate transaction creation failure during retries
	mockFaultTolerance.On("RetryOperation", mock.Anything, 3).Return(errors.New("transaction creation failure"))

	service := NewUserService(mockUserRepo, mockTransactionRepo, mockCountryRepo, mockGatewayRepo, mockFaultTolerance, nil)

	// When
	err := service.ExecuteTransaction(ctx, txData)

	// Then
	assert.Error(t, err)
	assert.Equal(t, domain.ErrTransactionNotCreated, err)
	mockUserRepo.AssertExpectations(t)
	mockCountryRepo.AssertExpectations(t)
	mockGatewayRepo.AssertExpectations(t)
	mockFaultTolerance.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}
