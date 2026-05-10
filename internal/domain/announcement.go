package domain

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID             uuid.UUID
	Title          string
	URL            string
	AvailableFrom  time.Time
	AvailableUntil *time.Time
}
