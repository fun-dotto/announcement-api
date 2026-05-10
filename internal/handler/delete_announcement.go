package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
	"github.com/google/uuid"
)

func (h *Handler) AnnouncementsV1Delete(ctx context.Context, request api.AnnouncementsV1DeleteRequestObject) (api.AnnouncementsV1DeleteResponseObject, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return api.AnnouncementsV1Delete404Response{}, nil
	}
	err = h.announcementService.DeleteAnnouncement(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.AnnouncementsV1Delete404Response{}, nil
		}
		return nil, err
	}

	return api.AnnouncementsV1Delete204Response{}, nil
}
