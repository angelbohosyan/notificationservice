package models

type Notification struct {
	Message   string
	Recipient string
	Channel   string
}

type NotificationId struct {
	Id string
}
