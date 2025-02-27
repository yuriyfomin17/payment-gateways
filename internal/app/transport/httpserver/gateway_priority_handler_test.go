package httpserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"payment-gateway/internal/app/transport/httpserver/mocks"
	"testing"
)

const gatewayPriorityEndpoint = "/gateway/priority"

func TestHttpServer_UpdateGatewayPriority(t *testing.T) {
	testCases := []struct {
		name               string
		requestBody        []byte
		mockServiceError   error
		expectedStatusCode int
		setupMock          bool
		expectedMessage    string
	}{
		{
			name:               "Successful gateway priority update",
			requestBody:        []byte(`{"gt_id": 1, "priority": "high"}`),
			mockServiceError:   nil,
			expectedStatusCode: http.StatusOK,
			setupMock:          true,
			expectedMessage:    "success",
		},
		{
			name:               "Invalid request body",
			requestBody:        []byte(`invalid json`),
			expectedStatusCode: http.StatusBadRequest,
			setupMock:          false,
			expectedMessage:    "unsupported-content-type",
		},
		{
			name:               "Invalid gateway priority value",
			requestBody:        []byte(`{"gt_id": 1, "priority": "invalid"}`),
			mockServiceError:   nil,
			expectedStatusCode: http.StatusBadRequest,
			setupMock:          false,
			expectedMessage:    "invalid-gateway-priority",
		},
		{
			name:               "Gateway service error during update",
			requestBody:        []byte(`{"gt_id": 1, "priority": "high"}`),
			mockServiceError:   errors.New("service error"),
			expectedStatusCode: http.StatusBadRequest,
			setupMock:          true,
			expectedMessage:    "gateway-priority-not-updated",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize mocks and HTTP server
			mockedGatewayService := mocks.NewGatewayService(t)
			httpServer := NewHttpServer(nil, mockedGatewayService, nil)

			// Setup mock based on test case
			if tc.setupMock {
				mockedGatewayService.On("UpdateGatewayPriority", mock.Anything, mock.Anything, mock.Anything).
					Return(tc.mockServiceError)
			}

			// Prepare request
			req := httptest.NewRequest(http.MethodPost, gatewayPriorityEndpoint, bytes.NewBuffer(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Call handler
			httpServer.UpdateGatewayPriority(w, req)

			// Capture response
			res := w.Result()
			defer res.Body.Close()

			// Validate response
			require.Equal(t, tc.expectedStatusCode, res.StatusCode)

			responseBody, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)

			var apiResponse APIResponse
			err = json.Unmarshal(responseBody, &apiResponse)
			require.NoError(t, err)

			require.Equal(t, tc.expectedMessage, apiResponse.Message)
			require.Equal(t, tc.expectedStatusCode, apiResponse.StatusCode)
		})
	}
}
