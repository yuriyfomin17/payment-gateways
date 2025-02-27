package httpserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"payment-gateway/internal/app/domain"
	"payment-gateway/internal/app/transport/httpserver/mocks"
	"testing"
)

const depositEndpoint = "/deposit"
const withdrawEndpoint = "/withdraw"

func TestHttpServer_DepositHandler(t *testing.T) {
	// given
	mockedTransactionPublisherService := mocks.NewTransactionPublisherService(t)
	mockedUserService := mocks.NewUserService(t)
	mockedGatewayService := mocks.NewGatewayService(t)
	requestBody := []byte(`{
	    "amount": 100.00,
	    "user_id": 1,
	    "currency": "EUR"
	}`)

	mockedUserService.On("FetchTxId", mock.Anything).Return(int64(1), nil)
	mockedUserService.On("ExecuteTransaction", mock.Anything, mock.Anything).Return(nil)

	httpServer := NewHttpServer(mockedUserService, mockedGatewayService, mockedTransactionPublisherService)

	req := httptest.NewRequest(http.MethodPost, depositEndpoint, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// when
	httpServer.DepositHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	// then
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Read response body
	responseBody, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	// Unmarshal response body
	var apiResponse APIResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, apiResponse.StatusCode)
	require.Equal(t, http.StatusOK, apiResponse.StatusCode)
	require.Equal(t, "success", apiResponse.Message)
	require.NotNil(t, apiResponse.Data)
	fmt.Println(apiResponse.Data)
	require.EqualValues(t, map[string]any{
		"transaction_id":     float64(1),
		"transaction_status": "pending",
	}, apiResponse.Data.(map[string]any))
}

func TestHttpServer_WithdrawalHandler(t *testing.T) {
	// given
	mockedTransactionPublisherService := mocks.NewTransactionPublisherService(t)
	mockedUserService := mocks.NewUserService(t)
	mockedGatewayService := mocks.NewGatewayService(t)
	requestBody := []byte(`{
	    "amount": 50.00, 
	    "user_id": 2,     
	    "currency": "USD"
	}`)

	mockedUserService.On("FetchTxId", mock.Anything).Return(int64(2), nil)
	mockedUserService.On("ExecuteTransaction", mock.Anything, mock.Anything).Return(nil)

	httpServer := NewHttpServer(mockedUserService, mockedGatewayService, mockedTransactionPublisherService)

	req := httptest.NewRequest(http.MethodPost, withdrawEndpoint, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// when
	httpServer.WithdrawalHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	// then
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Read response body
	responseBody, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	// Unmarshal response body
	var apiResponse APIResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, apiResponse.StatusCode)
	require.Equal(t, "success", apiResponse.Message)
	require.NotNil(t, apiResponse.Data)

	fmt.Println(apiResponse.Data)
	require.EqualValues(t, map[string]any{
		"transaction_id":     float64(2), // Verify the correct txID
		"transaction_status": "pending",
	}, apiResponse.Data.(map[string]any))
}

func TestHttpServer_ShouldHandleInvalidTransactionCreation(t *testing.T) {

	requestBody := []byte(`{"amount": 100.00, "user_id": 1, "currency": "EUR"}`)
	mockUserServiceError := errors.New("test error")
	expectedStatusCode := http.StatusInternalServerError
	expectedMessage := "transaction-not-created"

	mockedTransactionPublisherService := mocks.NewTransactionPublisherService(t)
	mockedUserService := mocks.NewUserService(t)
	mockedGatewayService := mocks.NewGatewayService(t)
	httpServer := NewHttpServer(mockedUserService, mockedGatewayService, mockedTransactionPublisherService)

	mockedUserService.On("FetchTxId", mock.Anything).Return(int64(0), mockUserServiceError)
	if mockUserServiceError == nil || errors.Is(mockUserServiceError, domain.ErrTransactionNotCreated) {
		mockedUserService.On("ExecuteTransaction", mock.Anything, mock.Anything).Return(mockUserServiceError)
	}

	req := httptest.NewRequest(http.MethodPost, depositEndpoint, bytes.NewBuffer(requestBody))

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	httpServer.DepositHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, expectedStatusCode, res.StatusCode)

	responseBody, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	var apiResponse APIResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	require.NoError(t, err)

	require.Equal(t, expectedStatusCode, apiResponse.StatusCode)

	require.Equal(t, expectedMessage, apiResponse.Message)
}

func TestHttpServer_ShouldHandleUnsupportedContentType(t *testing.T) {
	requestBody := []byte(`invalid json`)
	expectedStatusCode := http.StatusBadRequest
	expectedMessage := "unsupported-content-type"

	mockedTransactionPublisherService := mocks.NewTransactionPublisherService(t)
	mockedUserService := mocks.NewUserService(t)
	mockedGatewayService := mocks.NewGatewayService(t)
	httpServer := NewHttpServer(mockedUserService, mockedGatewayService, mockedTransactionPublisherService)

	req := httptest.NewRequest(http.MethodPost, depositEndpoint, bytes.NewBuffer(requestBody))

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	httpServer.DepositHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, expectedStatusCode, res.StatusCode)

	responseBody, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	var apiResponse APIResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	require.NoError(t, err)

	require.Equal(t, expectedStatusCode, apiResponse.StatusCode)

	require.Equal(t, expectedMessage, apiResponse.Message)
}

func TestHttpServer_ShouldHandleTransactionCreationError(t *testing.T) {
	requestBody := []byte(`{"amount": 100.00, "user_id": 1, "currency": "EUR"}`)
	expectedStatusCode := http.StatusBadRequest
	expectedMessage := "transaction-not-created"

	mockedTransactionPublisherService := mocks.NewTransactionPublisherService(t)
	mockedUserService := mocks.NewUserService(t)
	mockedGatewayService := mocks.NewGatewayService(t)
	httpServer := NewHttpServer(mockedUserService, mockedGatewayService, mockedTransactionPublisherService)

	mockedUserService.On("FetchTxId", mock.Anything).Return(int64(0), domain.ErrTransactionNotCreated)

	req := httptest.NewRequest(http.MethodPost, depositEndpoint, bytes.NewBuffer(requestBody))

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	httpServer.DepositHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, expectedStatusCode, res.StatusCode)

	responseBody, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	var apiResponse APIResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	require.NoError(t, err)

	require.Equal(t, expectedStatusCode, apiResponse.StatusCode)

	require.Equal(t, expectedMessage, apiResponse.Message)
}

func TestHttpServer_ShouldHandleErrorDuringTransactionExecution(t *testing.T) {
	requestBody := []byte(`{"amount": 100.00, "user_id": 1, "currency": "EUR"}`)
	mockUserServiceError := errors.New("test error")
	expectedStatusCode := http.StatusBadRequest
	expectedMessage := "transaction-not-created"

	mockedTransactionPublisherService := mocks.NewTransactionPublisherService(t)
	mockedUserService := mocks.NewUserService(t)
	mockedGatewayService := mocks.NewGatewayService(t)
	httpServer := NewHttpServer(mockedUserService, mockedGatewayService, mockedTransactionPublisherService)

	mockedUserService.On("FetchTxId", mock.Anything).Return(int64(2), nil)
	mockedUserService.On("ExecuteTransaction", mock.Anything, mock.Anything).Return(mockUserServiceError)

	req := httptest.NewRequest(http.MethodPost, depositEndpoint, bytes.NewBuffer(requestBody))

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	httpServer.DepositHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	require.Equal(t, expectedStatusCode, res.StatusCode)

	responseBody, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	var apiResponse APIResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	require.NoError(t, err)

	require.Equal(t, expectedStatusCode, apiResponse.StatusCode)

	require.Equal(t, expectedMessage, apiResponse.Message)
}
