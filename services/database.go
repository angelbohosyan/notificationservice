package services

import (
	"awesomeProject/models"
	"fmt"
	"log"
)

type NotificationResult struct {
	Notification *models.Notification
	Error        string
}

var notificationResult = make(map[string]*NotificationResult)

func addNotification(id string, notification *models.Notification, err error) {
	var errorValue string

	if err != nil {
		errorValue = err.Error()
	} else {
		errorValue = ""
	}

	notificationResult[id] = &NotificationResult{
		Notification: notification,
		Error:        errorValue,
	}
}

func pullNotification(id string) (*NotificationResult, error) {
	var val *NotificationResult
	var ok bool
	if val, ok = notificationResult[id]; !ok {
		log.Printf("can not find result with this id %v\n", id)
		return nil, fmt.Errorf("can not find result with this id %v", id)
	}

	return val, nil
}
