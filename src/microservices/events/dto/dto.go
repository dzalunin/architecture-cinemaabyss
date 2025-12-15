package dto

import (
	"time"
)

type MovieEvent struct {
	MovieID     int      `json:"movie_id" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	Action      string   `json:"action" binding:"required"`
	UserID      *int     `json:"user_id,omitempty"`
	Rating      *float64 `json:"rating,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Description *string  `json:"description,omitempty"`
}

type UserEvent struct {
	UserID    int       `json:"user_id" binding:"required"`
	Username  *string   `json:"username,omitempty"`
	Email     *string   `json:"email,omitempty"`
	Action    string    `json:"action" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required"`
}

type PaymentEvent struct {
	PaymentID  int       `json:"payment_id" binding:"required"`
	UserID     int       `json:"user_id" binding:"required"`
	Amount     float64   `json:"amount" binding:"required"`
	Status     string    `json:"status" binding:"required"`
	Timestamp  time.Time `json:"timestamp" binding:"required"`
	MethodType *string   `json:"method_type,omitempty"`
}

type EventResponse struct {
	Status    string `json:"status" binding:"required"`
	Partition int32  `json:"partition" binding:"required"`
	Offset    int64  `json:"offset" binding:"required"`
	Event     Event  `json:"event" binding:"required"`
}

type Event struct {
	ID        string    `json:"id" binding:"required"`
	Type      string    `json:"type" binding:"required"`
	Timestamp time.Time `json:"timestamp" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	Payload   any       `json:"payload" binding:"required"`
}
