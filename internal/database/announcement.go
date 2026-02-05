package database

import (
	"time"

	"github.com/fun-dotto/announcement-api/internal/domain"
)

type Announcement struct {
	ID             string     `gorm:"primaryKey;type:uuid"`
	Title          string     `gorm:"not null"`
	URL            string     `gorm:"not null"`
	AvailableFrom  time.Time  `gorm:"not null;index"`
	AvailableUntil *time.Time `gorm:"index"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (m *Announcement) ToDomain() domain.Announcement {
	return domain.Announcement{
		ID:             m.ID,
		Title:          m.Title,
		URL:            m.URL,
		AvailableFrom:  m.AvailableFrom,
		AvailableUntil: m.AvailableUntil,
	}
}

func FromDomain(announcement domain.Announcement) Announcement {
	return Announcement{
		ID:             announcement.ID,
		Title:          announcement.Title,
		URL:            announcement.URL,
		AvailableFrom:  announcement.AvailableFrom,
		AvailableUntil: announcement.AvailableUntil,
	}
}
