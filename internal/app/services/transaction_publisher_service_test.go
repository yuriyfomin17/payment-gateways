package services

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"payment-gateway/internal/app/domain"
	"payment-gateway/internal/app/services/mocks"
	"testing"
)

func TestPublishTransaction_JSON_Success(t *testing.T) {
	// Given
	mockFaultTolerance := mocks.NewFaultTolerance(t)
	mockRabbit := mocks.NewRabbitMqService(t)
	mockEncryptor := mocks.NewDataEncryptor(t)
	mockRedis := mocks.NewRedisService(t)

	service := TransactionPublisherServiceImpl{
		faultToleranceService:  mockFaultTolerance,
		rabbitService:          mockRabbit,
		dataEncryptor:          mockEncryptor,
		failedTransactionRedis: mockRedis,
	}

	ctx := context.TODO()
	txID := int64(1)
	statusToUpdate := "Success"
	dataFormat := "application/json"

	jsonMessage := domain.TxStatusRabbitJson{TxID: txID, Status: statusToUpdate}
	jsonBytes, _ := json.Marshal(jsonMessage)
	encryptedJsonBytes := []byte("encrypted-json")

	mockEncryptor.On("MaskData", jsonBytes).Return(encryptedJsonBytes, nil)
	mockFaultTolerance.On("PublishWithCircuitBreaker", mock.AnythingOfType("func() error")).Return(nil)

	// When
	err := service.PublishTransaction(ctx, txID, statusToUpdate, dataFormat)

	// Then
	assert.NoError(t, err)
	mockEncryptor.AssertCalled(t, "MaskData", jsonBytes)
}

func TestPublishTransaction_XML_Success(t *testing.T) {
	// Given
	mockFaultTolerance := mocks.NewFaultTolerance(t)
	mockRabbit := mocks.NewRabbitMqService(t)
	mockEncryptor := mocks.NewDataEncryptor(t)
	mockRedis := mocks.NewRedisService(t)

	service := TransactionPublisherServiceImpl{
		faultToleranceService:  mockFaultTolerance,
		rabbitService:          mockRabbit,
		dataEncryptor:          mockEncryptor,
		failedTransactionRedis: mockRedis,
	}

	ctx := context.TODO()
	txID := int64(2)
	statusToUpdate := "Failed"
	dataFormat := "application/xml"

	xmlMessage := domain.TxStatusXML{TxID: txID, Status: statusToUpdate}
	xmlBytes, _ := xml.Marshal(xmlMessage)
	encryptedXmlBytes := []byte("encrypted-xml")

	mockEncryptor.On("MaskData", xmlBytes).Return(encryptedXmlBytes, nil)
	mockFaultTolerance.On("PublishWithCircuitBreaker", mock.AnythingOfType("func() error")).Return(nil)

	// When
	err := service.PublishTransaction(ctx, txID, statusToUpdate, dataFormat)

	// Then
	assert.NoError(t, err)
	mockEncryptor.AssertCalled(t, "MaskData", xmlBytes)
}

func TestPublishTransaction_InvalidDataFormat(t *testing.T) {
	// Given
	mockFaultTolerance := mocks.NewFaultTolerance(t)
	mockRabbit := mocks.NewRabbitMqService(t)
	mockEncryptor := mocks.NewDataEncryptor(t)
	mockRedis := mocks.NewRedisService(t)

	service := TransactionPublisherServiceImpl{
		faultToleranceService:  mockFaultTolerance,
		rabbitService:          mockRabbit,
		dataEncryptor:          mockEncryptor,
		failedTransactionRedis: mockRedis,
	}

	ctx := context.TODO()
	txID := int64(4)
	statusToUpdate := "Completed"
	dataFormat := "text/plain"

	mockRedis.On("CreateFailedCallbackTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// When
	err := service.PublishTransaction(ctx, txID, statusToUpdate, dataFormat)

	// Then
	assert.Error(t, err)
	assert.Equal(t, domain.ErrTransactionStatus, err)
	mockEncryptor.AssertNotCalled(t, "MaskData", mock.Anything)
	mockFaultTolerance.AssertNotCalled(t, "PublishWithCircuitBreaker", mock.Anything)
	mockRedis.AssertCalled(t, "CreateFailedCallbackTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
