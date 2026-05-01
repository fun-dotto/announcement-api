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

func TestAnnouncementsV1Detail(t *testing.T) {
	id := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	t.Run("正常にお知らせを取得できる", func(t *testing.T) {
		var receivedID uuid.UUID
		mockRepo := &repository.MockAnnouncementRepository{
			GetAnnouncementByIDFunc: func(ctx context.Context, gotID uuid.UUID) (domain.Announcement, error) {
				receivedID = gotID
				return domain.Announcement{
					ID:            id,
					Title:         "詳細",
					URL:           "https://example.com/detail",
					AvailableFrom: now,
				}, nil
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		request := api.AnnouncementsV1DetailRequestObject{Id: id.String()}
		response, err := h.AnnouncementsV1Detail(context.Background(), request)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		assert.NoError(t, response.VisitAnnouncementsV1DetailResponse(w))
		assert.Equal(t, http.StatusOK, w.Code)

		var got api.AnnouncementsV1Detail200JSONResponse
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
		assert.Equal(t, id.String(), got.Announcement.Id)
		assert.Equal(t, "詳細", got.Announcement.Title)
		assert.Equal(t, id, receivedID)
	})

	t.Run("不正なUUIDで404を返す", func(t *testing.T) {
		called := false
		mockRepo := &repository.MockAnnouncementRepository{
			GetAnnouncementByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Announcement, error) {
				called = true
				return domain.Announcement{}, nil
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		request := api.AnnouncementsV1DetailRequestObject{Id: "not-a-uuid"}
		response, err := h.AnnouncementsV1Detail(context.Background(), request)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		assert.NoError(t, response.VisitAnnouncementsV1DetailResponse(w))
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.False(t, called, "repository must not be called for invalid UUID")
	})

	t.Run("リポジトリの一般エラーをそのまま返す", func(t *testing.T) {
		repoErr := errors.New("connection lost")
		mockRepo := &repository.MockAnnouncementRepository{
			GetAnnouncementByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Announcement, error) {
				return domain.Announcement{}, repoErr
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		request := api.AnnouncementsV1DetailRequestObject{Id: id.String()}
		response, err := h.AnnouncementsV1Detail(context.Background(), request)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, response)
	})
}

func TestAnnouncementsV1Detail_NotFound(t *testing.T) {
	id := uuid.MustParse("33333333-3333-3333-3333-333333333333")
	mockRepo := &repository.MockAnnouncementRepository{
		GetAnnouncementByIDFunc: func(ctx context.Context, id uuid.UUID) (domain.Announcement, error) {
			return domain.Announcement{}, domain.ErrNotFound
		},
	}
	h := NewHandler(service.NewAnnouncementService(mockRepo))

	request := api.AnnouncementsV1DetailRequestObject{Id: id.String()}
	response, err := h.AnnouncementsV1Detail(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	w := httptest.NewRecorder()
	err = response.VisitAnnouncementsV1DetailResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
