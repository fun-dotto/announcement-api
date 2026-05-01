package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
	"github.com/fun-dotto/announcement-api/internal/repository"
	"github.com/fun-dotto/announcement-api/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAnnouncementsV1Delete_NotFound(t *testing.T) {
	mockRepo := &repository.MockAnnouncementRepository{
		DeleteAnnouncementFunc: func(ctx context.Context, id uuid.UUID) error {
			return domain.ErrNotFound
		},
	}
	h := NewHandler(service.NewAnnouncementService(mockRepo))

	request := api.AnnouncementsV1DeleteRequestObject{Id: "nonexistent"}
	response, err := h.AnnouncementsV1Delete(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	w := httptest.NewRecorder()
	err = response.VisitAnnouncementsV1DeleteResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
