package telegram

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gelleson/changescout/changescout/internal/app/services/sender/providers/telegram/mocks"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTelegram_Send(t *testing.T) {
	mockDoer := mocks.NewDoer(t)                // Create a new Doer mock
	telegramClient := &Telegram{doer: mockDoer} // Instantiate the Telegram client with the mock

	tests := []struct {
		name          string
		notification  string
		conf          domain.Notification
		mockFunc      func()
		expectedError bool
	}{
		{
			name:         "successful send",
			notification: "Test notification",
			conf: domain.Notification{
				Token:       stringPtr("mock_token"),
				Destination: stringPtr("mock_chat_id"),
			},
			mockFunc: func() {
				// Mock a successful HTTP response
				response := httptest.NewRecorder()
				response.WriteHeader(http.StatusOK)
				mockDoer.On("Do", mock.AnythingOfType("*http.Request")).
					Return(response.Result(), nil).Once()
			},
			expectedError: false,
		},
		{
			name:         "failure in request execution",
			notification: "Test notification",
			conf: domain.Notification{
				Token:       stringPtr("mock_token"),
				Destination: stringPtr("mock_chat_id"),
			},
			mockFunc: func() {
				// Simulate a failure in request execution
				mockDoer.On("Do", mock.AnythingOfType("*http.Request")).
					Return(nil, errors.New("failed to execute request")).Once()
			},
			expectedError: true,
		},
		{
			name:         "non-200 status code",
			notification: "Test notification",
			conf: domain.Notification{
				Token:       stringPtr("mock_token"),
				Destination: stringPtr("mock_chat_id"),
			},
			mockFunc: func() {
				// Mock a non-200 HTTP response
				response := httptest.NewRecorder()
				response.WriteHeader(http.StatusBadRequest)
				mockDoer.On("Do", mock.AnythingOfType("*http.Request")).
					Return(response.Result(), nil).Once()
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc() // Setup the mock as per the test case

			// Run the method we want to test
			err := telegramClient.Send(tt.notification, tt.conf)

			// Assert expectations
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Validate all expectations for the mock were met, if applicable
			mockDoer.AssertExpectations(t)
		})
	}
}

// Helper function to get a pointer to a string
func stringPtr(s string) *string {
	return &s
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "default client",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			telegramClient := New()
			assert.NotNil(t, telegramClient, "Telegram instance should not be nil")
			assert.Equal(t, http.DefaultClient, telegramClient.doer, "Doer should be http.DefaultClient")
		})
	}
}

// Mocking json.Marshal function for testing purposes to control its behavior
var jsonMarshalFn = json.Marshal

// TestEncode tests the encode function.
func TestEncode(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectPanic bool
	}{
		{
			name:        "successful encoding",
			input:       map[string]string{"key": "value"},
			expectPanic: false,
		},
		{
			name:        "nil input encoding",
			input:       nil,
			expectPanic: false,
		},
		{
			name:        "encoding failure",
			input:       make(chan int),
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("The code did not panic for test case %s", tt.name)
					}
				}()
			}
			_ = encode(tt.input)
		})
	}
}
