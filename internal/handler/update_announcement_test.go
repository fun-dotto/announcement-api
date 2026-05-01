package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
	"github.com/fun-dotto/announcement-api/internal/repository"
	"github.com/fun-dotto/announcement-api/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestAnnouncementsV1Update_NotFound(t *testing.T) {
	mockRepo := &repository.MockAnnouncementRepository{
		UpdateAnnouncementFunc: func(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
			return domain.Announcement{}, domain.ErrNotFound
		},
	}
	h := NewHandler(service.NewAnnouncementService(mockRepo))

	body := api.AnnouncementsV1UpdateJSONRequestBody{
		Title:         "test",
		AvailableFrom: time.Now(),
		Url:           "https://example.com",
	}
	request := api.AnnouncementsV1UpdateRequestObject{Id: "nonexistent", Body: &body}
	response, err := h.AnnouncementsV1Update(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	w := httptest.NewRecorder()
	err = response.VisitAnnouncementsV1UpdateResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
