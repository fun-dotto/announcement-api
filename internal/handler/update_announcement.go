package handler

import (
	"context"
	"errors"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
	"github.com/google/uuid"
)

func (h *Handler) AnnouncementsV1Update(ctx context.Context, request api.AnnouncementsV1UpdateRequestObject) (api.AnnouncementsV1UpdateResponseObject, error) {
	id, err := uuid.Parse(request.Id)
	if err != nil {
		return api.AnnouncementsV1Update404Response{}, nil
	}
	domainAnnouncement := toDomainAnnouncementFromRequest(id, *request.Body)

	updated, err := h.announcementService.UpdateAnnouncement(ctx, domainAnnouncement)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return api.AnnouncementsV1Update404Response{}, nil
		}
		return nil, err
	}

	return api.AnnouncementsV1Update200JSONResponse{
		Announcement: toApiAnnouncement(updated),
	}, nil
}
