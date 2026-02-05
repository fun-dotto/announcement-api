package database

import (
	"time"

	"github.com/fun-dotto/announcement-api/internal/domain"
)

type Announcement struct {
	ID             string     `gorm:"primaryKey;type:uuid"`
	Title          string     `gorm:"type:varchar(500);not null"`
	URL            string     `gorm:"type:varchar(1000);not null"`
	AvailableFrom  *time.Time `gorm:"index"`
	AvailableUntil *time.Time `gorm:"index"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// v0廃止まで残す
	IsActive *bool `gorm:"default:true;index"`
	Date     *time.Time
}

func (Announcement) TableName() string {
	return "announcements"
}

func (m *Announcement) ToDomain() domain.Announcement {
	return domain.Announcement{
		ID:    m.ID,
		Title: m.Title,
		URL:   m.URL,
		AvailableFrom: func() time.Time {
			if m.AvailableFrom == nil {
				if m.Date == nil {
					return time.Now()
				}
				return *m.Date
			}
			return *m.AvailableFrom
		}(),
		AvailableUntil: m.AvailableUntil,
	}
}

func FromDomain(announcement domain.Announcement) Announcement {
	return Announcement{
		ID:             announcement.ID,
		Title:          announcement.Title,
		URL:            announcement.URL,
		AvailableFrom:  &announcement.AvailableFrom,
		AvailableUntil: announcement.AvailableUntil,
		Date:           &announcement.AvailableFrom,
		IsActive: func() *bool {
			b := announcement.AvailableUntil == nil || announcement.AvailableUntil.After(time.Now())
			return &b
		}(),
	}
}
