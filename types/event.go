package types

import "ems/models"

type EventCreateRequest struct {
	Title       string  `json:"title"`
	Location    *string `json:"location"`
	Description *string `json:"description"`
	StartTime   *string `json:"start_time"` // ISO8601 format recommended
	EndTime     *string `json:"end_time"`   // ISO8601 format recommended
	CreatedBy   *string `json:"created_by"`
}

type EventCreateResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Location    *string `json:"location"`
	Description *string `json:"description"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
	CreatedBy   *string `json:"created_by"`
	Message     string  `json:"message"`
}

// If you need the second struct, rename it, for example:
type EventCreateFullResponse struct {
	Message string        `json:"message"`
	Event   *models.Event `json:"event"`
}
type EventUpdateRequest struct{}
