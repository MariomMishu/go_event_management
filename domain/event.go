package domain

import (
	"ems/models"
)

// Event represents an event entity.
type Event struct {
	ID          int
	Title       string
	Description string
	StartTime   string // Consider using time.Time in real applications
	EndTime     string
	Location    string
}

// EventRepository defines the interface for event data operations.
type EventRepository interface {
	CreateEvent(event *models.Event) (*models.Event, error)
	ReadEventByID(id int) (*models.Event, error)
	UpdateEvent(event *models.Event) error
	DeleteEvent(id int) error
	ListEvents() ([]*models.Event, error)
}

// EventService defines the business logic for events.
type EventService interface {
	CreateEvent(event *models.Event) (*models.Event, error)
	GetEventByID(id int) (*models.Event, error)
	UpdateEvent(event *models.Event) error
	DeleteEvent(id int) error
	ListEvents() ([]*models.Event, error)
}
