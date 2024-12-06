package sender

import (
	"errors"
	"testing"

	"github.com/gelleson/changescout/changescout/internal/app/usecases/notification/mocks"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSenderService_Send(t *testing.T) {
	mockErr := errors.New("send error")
	notificationID := uuid.New()

	tests := []struct {
		name          string
		senders       Senders
		notification  string
		conf          domain.Notification
		expectedError error
	}{
		{
			name: "valid sender",
			senders: Senders{
				domain.NotificationType("email"): func() *mocks.Sender {
					mockSender := &mocks.Sender{}
					mockSender.On("Send", "This is a valid notification", mock.Anything).Return(nil)
					return mockSender
				}(),
			},
			notification: "This is a valid notification",
			conf: domain.Notification{
				ID:   notificationID,
				Type: domain.NotificationType("email"),
			},
			expectedError: nil,
		},
		{
			name: "unknown type",
			senders: Senders{
				domain.NotificationType("email"): &mocks.Sender{},
			},
			notification: "This is a valid notification",
			conf: domain.Notification{
				ID:   notificationID,
				Type: domain.NotificationType("sms"),
			},
			expectedError: nil,
		},
		{
			name: "sender returns error",
			senders: Senders{
				domain.NotificationType("email"): func() *mocks.Sender {
					mockSender := &mocks.Sender{}
					mockSender.On("Send", "This is a valid notification", mock.Anything).Return(mockErr)
					return mockSender
				}(),
			},
			notification: "This is a valid notification",
			conf: domain.Notification{
				ID:   notificationID,
				Type: domain.NotificationType("email"),
			},
			expectedError: mockErr,
		},
		{
			name: "empty notification",
			senders: Senders{
				domain.NotificationType("email"): func() *mocks.Sender {
					mockSender := &mocks.Sender{}
					mockSender.On("Send", "", mock.Anything).Return(nil)
					return mockSender
				}(),
			},
			notification: "",
			conf: domain.Notification{
				ID:   notificationID,
				Type: domain.NotificationType("email"),
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewSenderService(tt.senders)
			err := service.Send(tt.notification, tt.conf)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
