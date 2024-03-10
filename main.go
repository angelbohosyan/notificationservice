package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"notificationservice/models"
	"notificationservice/services"
)

func main() {
	http.HandleFunc("/message", message)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func message(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		sendMessage(w, r)
	}

	if r.Method == http.MethodGet {
		getMessage(w, r)
	}
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	var result []byte
	var err error

	query := r.URL.Query()
	id := query.Get("id")
	if id != "" {
		if result, err = json.Marshal(services.ApiError{Error: "no parameter error"}); err != nil {
			log.Printf("error marshaling: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	val, err := services.GetNotificationResult(id)
	if err != nil {
		if result, err = json.Marshal(services.ApiError{Error: "your notification is still processing or does not exist"}); err != nil {
			log.Printf("error marshaling: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	if result, err = json.Marshal(val.Notification); err != nil {
		log.Printf("could not marshal notification result: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if _, err = w.Write(result); err != nil {
		log.Printf("error returning body: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		log.Printf("you returned notification to the user: %s\n", err)
	}

}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	result, err := notificationResult(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var body []byte
	if body, err = json.Marshal(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if _, err = w.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	switch result.(type) {
	case services.ApiError:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func notificationResult(r *http.Request) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return services.ApiError{Error: "empty body"}, nil
	}

	var notification models.Notification
	if err = json.Unmarshal(body, &notification); err != nil {
		return nil, fmt.Errorf("could not unmarshal body: %s\n", err)
	}

	var nr *services.NotificationResult

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("error getting id %s", err)
	}

	if nr, err = services.SendNotification(&notification, id.String()); err != nil {
		log.Printf("could not send notification: %s\n", err)
		return services.ApiError{Error: err.Error()}, nil
	}

	if nr == nil {
		return services.ApiID{Id: id.String()}, nil
	} else if nr.Error != "" {
		return services.ApiError{Error: nr.Error}, nil
	}

	return nr.Notification, nil
}
