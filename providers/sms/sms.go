package sms

import (
	"awesomeProject/models"
	"awesomeProject/providers"
	"encoding/json"
	"github.com/sfreiberg/gotwilio"
	"log"
	"os"
	"time"
)

type Service struct {
	TwilioClient *gotwilio.Twilio
}

type SMSServiceConfiguration struct {
	sid   string
	token string
}

func NewSMSService() *Service {
	credsFile := "C:\\Users\\Angel\\Downloads\\smscreds.json"
	creds, err := os.ReadFile(credsFile)
	if err != nil {
		log.Fatalf("Error reading credentials file: %v", err)
	}

	var token SMSServiceConfiguration

	if err = json.Unmarshal(creds, &token); err != nil {
		log.Fatalf("can not get sms token")
	}

	twilio := gotwilio.NewTwilioClient(token.sid, token.token)
	return &Service{
		TwilioClient: twilio,
	}
}

func (s *Service) SendNotification(notification *models.Notification) error {
	return providers.Retry(3, time.Millisecond, func() error {
		_, _, err := s.TwilioClient.SendSMS("+12694212604", notification.Recipient, notification.Message, "", "")
		if err != nil {
			log.Println(err)
			return err
		}

		return nil
	})
}
