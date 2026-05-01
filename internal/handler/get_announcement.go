package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
	"github.com/google/uuid"
)

func (h *Handler) AnnouncementsV1Detail(ctx context.Context, request api.AnnouncementsV1DetailRequestObject) (api.AnnouncementsV1DetailResponseObject, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return api.AnnouncementsV1Detail404Response{}, nil
	}
	announcement, err := h.announcementService.GetAnnouncementByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.AnnouncementsV1Detail404Response{}, nil
		}
		return nil, err
	}

	return api.AnnouncementsV1Detail200JSONResponse{
		Announcement: toApiAnnouncement(announcement),
	}, nil
}
