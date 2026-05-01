package repository

import (
	"context"
	"errors"

	"github.com/fun-dotto/announcement-api/internal/domain"
	"github.com/fun-dotto/shared-go/db/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type announcementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) *announcementRepository {
	return &announcementRepository{db: db}
}

func (r *announcementRepository) GetAnnouncements(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	var dbAnnouncements []model.Announcement
	dbQuery := r.db.WithContext(ctx)

	if query.FilterIsActive {
		dbQuery = dbQuery.Where("available_from <= NOW()").Where("available_until IS NULL OR available_until > NOW()")
	}

	sortDateDirection := func() string {
		switch query.SortByDate {
		case domain.SortDirectionAsc:
			return "ASC"
		case domain.SortDirectionDesc:
			return "DESC"
		default:
			return "ASC"
		}
	}()

	dbQuery = dbQuery.Order("available_from " + sortDateDirection)

	if err := dbQuery.Find(&dbAnnouncements).Error; err != nil {
		return nil, err
	}

	domainAnnouncements := make([]domain.Announcement, len(dbAnnouncements))
	for i, dbAnnouncement := range dbAnnouncements {
		domainAnnouncements[i] = toDomainAnnouncement(dbAnnouncement)
	}

	return domainAnnouncements, nil
}

func (r *announcementRepository) GetAnnouncementByID(ctx context.Context, id uuid.UUID) (domain.Announcement, error) {
	var dbAnnouncement model.Announcement
	if err := r.db.WithContext(ctx).Where("id = ?", id.String()).First(&dbAnnouncement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Announcement{}, domain.ErrNotFound
		}
		return domain.Announcement{}, err
	}
	return toDomainAnnouncement(dbAnnouncement), nil
}

func (r *announcementRepository) CreateAnnouncement(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
	dbAnnouncement := toDatabaseAnnouncement(announcement)
	if err := r.db.WithContext(ctx).Create(&dbAnnouncement).Error; err != nil {
		return domain.Announcement{}, err
	}
	return toDomainAnnouncement(dbAnnouncement), nil
}

func (r *announcementRepository) UpdateAnnouncement(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
	dbAnnouncement := toDatabaseAnnouncement(announcement)
	result := r.db.WithContext(ctx).Model(&model.Announcement{}).Where("id = ?", announcement.ID).Updates(map[string]interface{}{
		"title":           dbAnnouncement.Title,
		"url":             dbAnnouncement.URL,
		"available_from":  dbAnnouncement.AvailableFrom,
		"available_until": dbAnnouncement.AvailableUntil,
	})
	if result.Error != nil {
		return domain.Announcement{}, result.Error
	}
	if result.RowsAffected == 0 {
		return domain.Announcement{}, domain.ErrNotFound
	}
	return r.GetAnnouncementByID(ctx, announcement.ID)
}

func (r *announcementRepository) DeleteAnnouncement(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Announcement{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
