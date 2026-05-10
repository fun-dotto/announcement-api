package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
	"github.com/fun-dotto/announcement-api/internal/repository"
	"github.com/fun-dotto/announcement-api/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAnnouncementsV1Update(t *testing.T) {
	id := uuid.MustParse("44444444-4444-4444-4444-444444444444")
	availableFrom := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	availableUntil := availableFrom.Add(24 * time.Hour)

	t.Run("正常にお知らせを更新できる", func(t *testing.T) {
		var captured domain.Announcement
		mockRepo := &repository.MockAnnouncementRepository{
			UpdateAnnouncementFunc: func(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
				captured = announcement
				return announcement, nil
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		body := api.AnnouncementsV1UpdateJSONRequestBody{
			Title:          "更新後タイトル",
			AvailableFrom:  availableFrom,
			AvailableUntil: &availableUntil,
			Url:            "https://example.com/updated",
		}
		request := api.AnnouncementsV1UpdateRequestObject{Id: id.String(), Body: &body}
		response, err := h.AnnouncementsV1Update(context.Background(), request)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		assert.NoError(t, response.VisitAnnouncementsV1UpdateResponse(w))
		assert.Equal(t, http.StatusOK, w.Code)

		var got api.AnnouncementsV1Update200JSONResponse
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
		assert.Equal(t, id.String(), got.Announcement.Id)
		assert.Equal(t, "更新後タイトル", got.Announcement.Title)
		assert.Equal(t, "https://example.com/updated", got.Announcement.Url)

		assert.Equal(t, id, captured.ID, "URL の id がドメインに反映される")
		assert.Equal(t, "更新後タイトル", captured.Title)
		assert.Equal(t, availableFrom, captured.AvailableFrom)
		assert.Equal(t, &availableUntil, captured.AvailableUntil)
		assert.Equal(t, "https://example.com/updated", captured.URL)
	})

	t.Run("不正なUUIDで404を返す", func(t *testing.T) {
		called := false
		mockRepo := &repository.MockAnnouncementRepository{
			UpdateAnnouncementFunc: func(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
				called = true
				return announcement, nil
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		body := api.AnnouncementsV1UpdateJSONRequestBody{
			Title:         "test",
			AvailableFrom: availableFrom,
			Url:           "https://example.com",
		}
		request := api.AnnouncementsV1UpdateRequestObject{Id: "not-a-uuid", Body: &body}
		response, err := h.AnnouncementsV1Update(context.Background(), request)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		assert.NoError(t, response.VisitAnnouncementsV1UpdateResponse(w))
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.False(t, called, "repository must not be called for invalid UUID")
	})

	t.Run("リポジトリの一般エラーをそのまま返す", func(t *testing.T) {
		repoErr := errors.New("write failed")
		mockRepo := &repository.MockAnnouncementRepository{
			UpdateAnnouncementFunc: func(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
				return domain.Announcement{}, repoErr
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		body := api.AnnouncementsV1UpdateJSONRequestBody{
			Title:         "test",
			AvailableFrom: availableFrom,
			Url:           "https://example.com",
		}
		request := api.AnnouncementsV1UpdateRequestObject{Id: id.String(), Body: &body}
		response, err := h.AnnouncementsV1Update(context.Background(), request)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, response)
	})
}

func TestAnnouncementsV1Update_NotFound(t *testing.T) {
	id := uuid.MustParse("55555555-5555-5555-5555-555555555555")
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
	request := api.AnnouncementsV1UpdateRequestObject{Id: id.String(), Body: &body}
	response, err := h.AnnouncementsV1Update(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	w := httptest.NewRecorder()
	err = response.VisitAnnouncementsV1UpdateResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
