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

func TestAnnouncementsV1Create(t *testing.T) {
	availableFrom := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	availableUntil := availableFrom.Add(72 * time.Hour)

	t.Run("正常にお知らせが作成できる", func(t *testing.T) {
		var captured domain.Announcement
		mockRepo := &repository.MockAnnouncementRepository{
			CreateAnnouncementFunc: func(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
				captured = announcement
				return announcement, nil
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		body := api.AnnouncementsV1CreateJSONRequestBody{
			Title:          "新しいお知らせ",
			AvailableFrom:  availableFrom,
			AvailableUntil: &availableUntil,
			Url:            "https://example.com/new",
		}
		request := api.AnnouncementsV1CreateRequestObject{Body: &body}
		response, err := h.AnnouncementsV1Create(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, response)

		w := httptest.NewRecorder()
		err = response.VisitAnnouncementsV1CreateResponse(w)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)

		var got api.AnnouncementsV1Create201JSONResponse
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
		assert.Equal(t, "新しいお知らせ", got.Announcement.Title)
		assert.Equal(t, "https://example.com/new", got.Announcement.Url)
		assert.Equal(t, availableFrom, got.Announcement.AvailableFrom)
		assert.NotNil(t, got.Announcement.AvailableUntil)
		assert.Equal(t, availableUntil, *got.Announcement.AvailableUntil)

		assert.NotEqual(t, uuid.Nil, captured.ID, "ID should be generated")
		assert.Equal(t, "新しいお知らせ", captured.Title)
		assert.Equal(t, "https://example.com/new", captured.URL)
		assert.Equal(t, availableFrom, captured.AvailableFrom)
		assert.Equal(t, &availableUntil, captured.AvailableUntil)
	})

	t.Run("リポジトリエラーをそのまま返す", func(t *testing.T) {
		repoErr := errors.New("db down")
		mockRepo := &repository.MockAnnouncementRepository{
			CreateAnnouncementFunc: func(ctx context.Context, announcement domain.Announcement) (domain.Announcement, error) {
				return domain.Announcement{}, repoErr
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		body := api.AnnouncementsV1CreateJSONRequestBody{
			Title:         "test",
			AvailableFrom: availableFrom,
			Url:           "https://example.com",
		}
		request := api.AnnouncementsV1CreateRequestObject{Body: &body}
		response, err := h.AnnouncementsV1Create(context.Background(), request)

		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, response)
	})
}
