package services

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"payment-gateway/internal/app/services/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"payment-gateway/internal/app/domain"
)

func TestCallbackService_StartListeningJsonMessages(t *testing.T) {
	// Mock dependencies
	mockRabbit := mocks.NewRabbitMqService(t)
	mockEncryptor := mocks.NewDataEncryptor(t)
	mockTxRepo := mocks.NewTransactionRepository(t)

	// Create the service instance
	callbackService := NewCallbackService(mockRabbit, mockEncryptor, mockTxRepo)

	// Set up a fake JSON message channel
	jsonMessageChan := make(chan []byte, 1)
	mockRabbit.On("GetJsonMessage").Return(jsonMessageChan)

	// Prepare a mocked JSON message
	mockTxStatus := domain.TxStatusRabbitJson{TxID: 12345, Status: "Completed"}
	mockTxJSON, _ := json.Marshal(mockTxStatus)
	jsonMessageChan <- mockTxJSON

	// Mock the unmasking process and repository behavior
	mockEncryptor.On("UnmaskData", mockTxJSON).Return(mockTxJSON, nil)
	mockTxRepo.On("UpdateTransactionStatus", mock.Anything, mockTxStatus.TxID, mockTxStatus.Status).Return(nil)

	// Start listening for messages
	go callbackService.StartListeningJsonMessages()

	// Wait for goroutine processing
	time.Sleep(100 * time.Millisecond)

	// Assertions
	mockRabbit.AssertExpectations(t)
	mockEncryptor.AssertExpectations(t)
	mockTxRepo.AssertExpectations(t)

	assert.NoError(t, nil, "No error should have occurred while processing JSON messages")
}

func TestCallbackService_StartListeningJsonMessages_ErrorUnmasking(t *testing.T) {
	// Mock dependencies
	mockRabbit := mocks.NewRabbitMqService(t)
	mockEncryptor := mocks.NewDataEncryptor(t)
	mockTxRepo := mocks.NewTransactionRepository(t)

	// Create the service instance
	callbackService := NewCallbackService(mockRabbit, mockEncryptor, mockTxRepo)

	// Set up a fake JSON message channel
	jsonMessageChan := make(chan []byte, 1)
	mockRabbit.On("GetJsonMessage").Return(jsonMessageChan)

	// Prepare a mocked JSON message
	mockTxJSON := []byte(`{invalid-json`)
	jsonMessageChan <- mockTxJSON

	// Mock the unmasking error
	mockEncryptor.On("UnmaskData", mockTxJSON).Return(nil, errors.New("unmasking error"))

	// Start listening for messages
	go callbackService.StartListeningJsonMessages()

	// Wait for goroutine processing
	time.Sleep(100 * time.Millisecond)

	// Assertions
	mockRabbit.AssertExpectations(t)
	mockEncryptor.AssertExpectations(t)
	mockTxRepo.AssertNotCalled(t, "UpdateTransactionStatus", mock.Anything, mock.Anything, mock.Anything)
}

func TestCallbackService_StartListeningSoapMessages(t *testing.T) {
	// Mock dependencies
	mockRabbit := mocks.NewRabbitMqService(t)
	mockEncryptor := mocks.NewDataEncryptor(t)
	mockTxRepo := mocks.NewTransactionRepository(t)

	// Create the service instance
	callbackService := NewCallbackService(mockRabbit, mockEncryptor, mockTxRepo)

	// Set up a fake SOAP message channel
	soapMessageChan := make(chan []byte, 1)
	mockRabbit.On("GetSoapMessage").Return(soapMessageChan)

	// Prepare a mocked SOAP message
	mockTxStatus := domain.TxStatusXML{TxID: 67890, Status: "Failed"}
	mockSoapXML, _ := xml.Marshal(mockTxStatus)
	soapMessageChan <- mockSoapXML

	// Mock the unmasking process and repository behavior
	mockEncryptor.On("UnmaskData", mockSoapXML).Return(mockSoapXML, nil)
	mockTxRepo.On("UpdateTransactionStatus", mock.Anything, mockTxStatus.TxID, mockTxStatus.Status).Return(nil)

	// Start listening for messages
	go callbackService.StartListeningSoapMessages()

	// Wait for goroutine processing
	time.Sleep(100 * time.Millisecond)

	// Assertions
	mockRabbit.AssertExpectations(t)
	mockEncryptor.AssertExpectations(t)
	mockTxRepo.AssertExpectations(t)

	assert.NoError(t, nil, "No error should have occurred while processing SOAP messages")
}

func TestCallbackService_StartListeningSoapMessages_ErrorUnmasking(t *testing.T) {
	// Mock dependencies
	mockRabbit := mocks.NewRabbitMqService(t)
	mockEncryptor := mocks.NewDataEncryptor(t)
	mockTxRepo := mocks.NewTransactionRepository(t)

	// Create the service instance
	callbackService := NewCallbackService(mockRabbit, mockEncryptor, mockTxRepo)

	// Set up a fake SOAP message channel
	soapMessageChan := make(chan []byte, 1)
	mockRabbit.On("GetSoapMessage").Return(soapMessageChan)

	// Prepare a mocked SOAP message
	mockSoapXML := []byte(`<invalid-xml`)
	soapMessageChan <- mockSoapXML

	// Mock the unmasking error
	mockEncryptor.On("UnmaskData", mockSoapXML).Return(nil, errors.New("unmasking error"))

	// Start listening for messages
	go callbackService.StartListeningSoapMessages()

	// Wait for goroutine processing
	time.Sleep(100 * time.Millisecond)

	// Assertions
	mockRabbit.AssertExpectations(t)
	mockEncryptor.AssertExpectations(t)
	mockTxRepo.AssertNotCalled(t, "UpdateTransactionStatus", mock.Anything, mock.Anything, mock.Anything)
}
