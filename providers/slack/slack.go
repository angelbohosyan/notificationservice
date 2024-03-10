package slack

import (
	"awesomeProject/models"
	"awesomeProject/providers"
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"log"
	"os"
	"time"
)

type Service struct {
	SlackClient *slack.Client
}

type SlackToken struct {
	Token string `json:"token"`
}

func NewSlackService() *Service {
	credsFile := "C:\\Users\\Angel\\Downloads\\slackcreds.json"
	creds, err := os.ReadFile(credsFile)
	if err != nil {
		log.Fatalf("Error reading credentials file: %v", err)
	}

	var token SlackToken

	if err = json.Unmarshal(creds, &token); err != nil {
		log.Fatalf("can not get slack token")
	}

	slackClient := slack.New(token.Token)
	return &Service{
		SlackClient: slackClient,
	}
}

func (s *Service) SendNotification(notification *models.Notification) error {
	userID, err := s.getUserID(notification.Recipient)
	if err != nil {
		return err
	}

	return providers.Retry(3, time.Millisecond, func() error {
		_, _, err = s.SlackClient.PostMessage(userID, slack.MsgOptionText(notification.Message, false))
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) getUserID(name string) (string, error) {
	users, err := s.SlackClient.GetUsers()
	if err != nil {
		return "", fmt.Errorf("error retrieving user list: %v", err)
	}

	for _, user := range users {
		if user.Name == name {
			return user.ID, nil
		}
	}

	return "", fmt.Errorf("did not find the specified user: %v", name)
}
