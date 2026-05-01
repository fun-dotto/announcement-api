package handler

import (
	"context"
	"errors"
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

func TestAnnouncementsV1Delete(t *testing.T) {
	id := uuid.MustParse("66666666-6666-6666-6666-666666666666")

	t.Run("正常にお知らせを削除できる", func(t *testing.T) {
		var received uuid.UUID
		mockRepo := &repository.MockAnnouncementRepository{
			DeleteAnnouncementFunc: func(ctx context.Context, gotID uuid.UUID) error {
				received = gotID
				return nil
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		request := api.AnnouncementsV1DeleteRequestObject{Id: id.String()}
		response, err := h.AnnouncementsV1Delete(context.Background(), request)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		assert.NoError(t, response.VisitAnnouncementsV1DeleteResponse(w))
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, id, received)
	})

	t.Run("不正なUUIDで404を返す", func(t *testing.T) {
		called := false
		mockRepo := &repository.MockAnnouncementRepository{
			DeleteAnnouncementFunc: func(ctx context.Context, id uuid.UUID) error {
				called = true
				return nil
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		request := api.AnnouncementsV1DeleteRequestObject{Id: "not-a-uuid"}
		response, err := h.AnnouncementsV1Delete(context.Background(), request)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		assert.NoError(t, response.VisitAnnouncementsV1DeleteResponse(w))
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.False(t, called, "repository must not be called for invalid UUID")
	})

	t.Run("リポジトリの一般エラーをそのまま返す", func(t *testing.T) {
		repoErr := errors.New("delete failed")
		mockRepo := &repository.MockAnnouncementRepository{
			DeleteAnnouncementFunc: func(ctx context.Context, id uuid.UUID) error {
				return repoErr
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		request := api.AnnouncementsV1DeleteRequestObject{Id: id.String()}
		response, err := h.AnnouncementsV1Delete(context.Background(), request)
		assert.ErrorIs(t, err, repoErr)
		assert.Nil(t, response)
	})
}

func TestAnnouncementsV1Delete_NotFound(t *testing.T) {
	id := uuid.MustParse("77777777-7777-7777-7777-777777777777")
	mockRepo := &repository.MockAnnouncementRepository{
		DeleteAnnouncementFunc: func(ctx context.Context, id uuid.UUID) error {
			return domain.ErrNotFound
		},
	}
	h := NewHandler(service.NewAnnouncementService(mockRepo))

	request := api.AnnouncementsV1DeleteRequestObject{Id: id.String()}
	response, err := h.AnnouncementsV1Delete(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	w := httptest.NewRecorder()
	err = response.VisitAnnouncementsV1DeleteResponse(w)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
