package domain

import "time"

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" validate:"required,min=5,max=100"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" validate:"required,gt"`
	EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
