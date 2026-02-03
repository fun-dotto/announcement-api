package database

import (
	"time"

	"github.com/fun-dotto/announcement-api/internal/domain"
)

type Announcement struct {
	ID             string     `gorm:"primaryKey;type:uuid"`
	Title          string     `gorm:"type:varchar(500);not null"`
	URL            string     `gorm:"type:varchar(1000);not null"`
	AvailableFrom  time.Time  `gorm:"index"`
	AvailableUntil *time.Time `gorm:"index"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// v0廃止まで残す
	IsActive bool      `gorm:"not null;default:true;index"`
	Date     time.Time `gorm:"not null;index"`
}

func (Announcement) TableName() string {
	return "announcements"
}

func (m *Announcement) ToDomain() domain.Announcement {
	return domain.Announcement{
		ID:       m.ID,
		Title:    m.Title,
		Date:     m.Date,
		URL:      m.URL,
		IsActive: m.IsActive,
	}
}

func FromDomain(announcement domain.Announcement) Announcement {
	return Announcement{
		ID:       announcement.ID,
		Title:    announcement.Title,
		Date:     announcement.Date,
		URL:      announcement.URL,
		IsActive: announcement.IsActive,
	}
}
