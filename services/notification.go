package services

import (
	"awesomeProject/models"
	"awesomeProject/providers/email"
	"awesomeProject/providers/slack"
	"awesomeProject/providers/sms"
	"fmt"
	"time"
)

type NotificationManager struct {
	EmailChannel *email.Service
	SMSChannel   *sms.Service
	SlackChannel *slack.Service
}

var notificationManager NotificationManager

func init() {
	notificationManager = NotificationManager{
		EmailChannel: email.NewEmailService(),
		SMSChannel:   sms.NewSMSService(),
		SlackChannel: slack.NewSlackService(),
	}
}

type ApiID struct {
	Id string
}

type ApiError struct {
	Error string
}

func SendNotification(notification *models.Notification, id string) (*NotificationResult, error) {
	switch notification.Channel {
	case "sms":
		go func() {
			notificationResultErr := notificationManager.SMSChannel.SendNotification(notification)
			addNotification(id, notification, notificationResultErr)
		}()
	case "email":
		go func() {
			notificationResultErr := notificationManager.EmailChannel.SendNotification(notification)
			addNotification(id, notification, notificationResultErr)
		}()
	case "slack":
		go func() {
			notificationResultErr := notificationManager.SlackChannel.SendNotification(notification)
			addNotification(id, notification, notificationResultErr)
		}()
	default:
		return nil, fmt.Errorf("unsupported channel: %s", notification.Channel)
	}

	time.Sleep(time.Second)

	val, err := pullNotification(id)
	if err != nil {
		return nil, nil
	}

	return val, nil
}

func GetNotificationResult(id string) (*NotificationResult, error) {
	return pullNotification(id)
}
