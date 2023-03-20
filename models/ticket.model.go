package models

import (
	"database/sql/driver"
	"time"
)

type ticketStatus string

const (
	NEW         ticketStatus = "NEW"
	IN_PROGRESS ticketStatus = "IN_PROGRESS"
	RESOLVED    ticketStatus = "RESOLVED"
)

func (ct *ticketStatus) Scan(value interface{}) error {
	*ct = ticketStatus(value.([]byte))
	return nil
}

func (ct ticketStatus) Value() (driver.Value, error) {
	return string(ct), nil
}

type Ticket struct {
	ID          uint      `gorm:"primary_key" json:"id,omitempty"`
	Title       string    `gorm:"not null" json:"title,omitempty"`
	Description string    `gorm:"not null" json:"description,omitempty"`
	Status      string    `gorm:"not null; type:ticket_status" json:"status,omitempty"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateTicketRequest struct {
	Title       string    `json:"title"  binding:"required"`
	Description string    `json:"description" binding:"required"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type UpdateTicket struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreateAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
