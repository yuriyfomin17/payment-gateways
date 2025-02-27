package services

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"payment-gateway/internal/app/services/mocks"
	"testing"
)

func TestGatewayService_UpdateGatewayPriority(t *testing.T) {
	// Arrange
	mockRepo := mocks.NewGatewayRepository(t)

	service := NewGatewayService(mockRepo)
	ctx := context.Background()
	gatewayID := int64(12345)
	priority := "High"

	// Mocking the behavior of UpdateGatewayPriority
	mockRepo.On("UpdateGatewayPriority", ctx, gatewayID, priority).Return(nil)

	// When
	err := service.UpdateGatewayPriority(ctx, gatewayID, priority)

	// Then
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGatewayService_UpdateGatewayPriority_Error(t *testing.T) {
	// Given
	mockRepo := mocks.NewGatewayRepository(t)

	service := NewGatewayService(mockRepo)
	ctx := context.Background()
	gatewayID := int64(12345)
	priority := "High"

	// Mocking the behavior of UpdateGatewayPriority to return an error
	mockRepo.On("UpdateGatewayPriority", ctx, gatewayID, priority).Return(errors.New("failed to update priority"))

	// When
	err := service.UpdateGatewayPriority(ctx, gatewayID, priority)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "failed to update priority", err.Error())
	mockRepo.AssertExpectations(t)
}
