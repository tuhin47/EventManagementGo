package router

import (
	"EventManagement/handler"
	"EventManagement/utils"
	"fmt"
	"log"
	"net/http"
)

func RegisterRoutes() {
	log.Printf("Registering HTTP handlers")
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/events", handleEvents)
	http.HandleFunc("/event/", handleEvent)

}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	utils.SetJSONHeader(w)
	switch r.Method {
	case http.MethodGet:
		handler.GetEventByIDHandler(w, r)
	case http.MethodPut:
		handler.UpdateEventHandler(w, r)
	case http.MethodDelete:
		handler.DeleteEventHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	utils.SetJSONHeader(w)
	switch r.Method {
	case http.MethodGet:
		handler.GetAllEventsHandler(w, r)
	case http.MethodPost:
		handler.CreateEventHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	utils.SetJSONHeader(w)
	fmt.Fprintf(w, `{"message": "Welcome to the Event Management System!"}`)
}
