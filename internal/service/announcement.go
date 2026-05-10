package service

import (
	"context"

	"github.com/fun-dotto/announcement-api/internal/domain"
	"github.com/google/uuid"
)

type AnnouncementRepository interface {
	GetAnnouncements(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error)
	GetAnnouncementByID(ctx context.Context, id uuid.UUID) (domain.Announcement, error)
	CreateAnnouncement(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error)
	UpdateAnnouncement(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error)
	DeleteAnnouncement(ctx context.Context, id uuid.UUID) error
}

type AnnouncementService struct {
	announcementRepository AnnouncementRepository
}

func NewAnnouncementService(announcementRepository AnnouncementRepository) *AnnouncementService {
	return &AnnouncementService{announcementRepository: announcementRepository}
}

func (s *AnnouncementService) GetAnnouncements(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
	return s.announcementRepository.GetAnnouncements(ctx, query)
}

func (s *AnnouncementService) GetAnnouncementByID(ctx context.Context, id uuid.UUID) (domain.Announcement, error) {
	return s.announcementRepository.GetAnnouncementByID(ctx, id)
}

func (s *AnnouncementService) CreateAnnouncement(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
	return s.announcementRepository.CreateAnnouncement(ctx, announcement)
}

func (s *AnnouncementService) UpdateAnnouncement(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
	return s.announcementRepository.UpdateAnnouncement(ctx, announcement)
}

func (s *AnnouncementService) DeleteAnnouncement(ctx context.Context, id uuid.UUID) error {
	return s.announcementRepository.DeleteAnnouncement(ctx, id)
}
