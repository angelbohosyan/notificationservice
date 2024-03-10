package email

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"notificationservice/models"
	"notificationservice/providers"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

type Service struct {
	client *http.Client
}

func NewEmailService() *Service {
	ctx := context.Background()

	credsFile := "C:\\Users\\Angel\\Downloads\\emailcreds.json"
	creds, err := os.ReadFile(credsFile)
	if err != nil {
		log.Fatalf("Error reading credentials file: %v", err)
	}

	config, err := google.ConfigFromJSON(creds, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Error parsing client secret file to config: %v", err)
	}

	token, err := getTokenFromUser(ctx, config)
	if err != nil {
		log.Fatalf("Error getting token from user: %v", err)
	}

	return &Service{client: config.Client(ctx, token)}
}

func (e *Service) SendNotification(notification *models.Notification) error {
	message := createMessage("angelbohosyan@gmail.com", notification.Recipient, "Subject", notification.Message)
	return providers.Retry(3, time.Millisecond, func() error {
		if err := e.sendEmail(message); err != nil {
			log.Println(err)
			return err
		}

		return nil
	})
}

func (e *Service) sendEmail(message gmail.Message) error {
	srv, err := gmail.New(e.client)
	if err != nil {
		return fmt.Errorf("unable to create Gmail service: %v", err)
	}

	if _, err = srv.Users.Messages.Send("me", &message).Do(); err != nil {
		return fmt.Errorf("unable to send email: %v", err)
	}

	log.Println("Email sent successfully!")

	return nil
}

func getTokenFromUser(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("authenticate:\n%v\n", authURL)
	log.Println("enter the authorization code: ")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}

	return token, nil
}

func createMessage(sender, recipient, subject, body string) gmail.Message {
	return gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", sender, recipient, subject, body))),
	}
}
