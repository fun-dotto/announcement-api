package repository

import (
	"github.com/fun-dotto/announcement-api/internal/domain"
	"github.com/fun-dotto/shared-go/db/model"
)

func toDatabaseAnnouncement(announcement domain.Announcement) model.Announcement {
	return model.Announcement{
		Common: model.Common{
			ID: announcement.ID,
		},
		Title:          announcement.Title,
		URL:            announcement.URL,
		AvailableFrom:  announcement.AvailableFrom,
		AvailableUntil: announcement.AvailableUntil,
	}
}

func toDomainAnnouncement(m model.Announcement) domain.Announcement {
	return domain.Announcement{
		ID:             m.ID,
		Title:          m.Title,
		URL:            m.URL,
		AvailableFrom:  m.AvailableFrom,
		AvailableUntil: m.AvailableUntil,
	}
}
