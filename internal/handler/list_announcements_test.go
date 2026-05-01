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

func TestAnnouncementsV1List(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	yesterday := now.Add(-24 * time.Hour)
	twoDaysAgo := now.Add(-48 * time.Hour)

	tests := []struct {
		name      string
		setupMock func() *repository.MockAnnouncementRepository
		params    api.AnnouncementsV1ListParams
		wantCode  int
		validate  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "正常にお知らせ一覧が取得できる",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{
					GetAnnouncementsFunc: func(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
						return []domain.Announcement{
							{ID: uuid.New(), Title: "お知らせ1", AvailableFrom: now, URL: "https://example.com/1"},
							{ID: uuid.New(), Title: "お知らせ2", AvailableFrom: yesterday, URL: "https://example.com/2"},
							{ID: uuid.New(), Title: "お知らせ3", AvailableFrom: twoDaysAgo, URL: "https://example.com/3"},
						}, nil
					},
				}
			},
			params:   api.AnnouncementsV1ListParams{},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response api.AnnouncementsV1List200JSONResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "failed to unmarshal response body")
				assert.Len(t, response.Announcements, 3)
			},
		},
		{
			name: "空の結果を正常に返せる",
			setupMock: func() *repository.MockAnnouncementRepository {
				return &repository.MockAnnouncementRepository{
					GetAnnouncementsFunc: func(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
						return []domain.Announcement{}, nil
					},
				}
			},
			params:   api.AnnouncementsV1ListParams{},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response api.AnnouncementsV1List200JSONResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Empty(t, response.Announcements)
			},
		},
		{
			name: "Announcementのフィールドがレスポンスに正しく変換される",
			setupMock: func() *repository.MockAnnouncementRepository {
				id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
				until := now.Add(48 * time.Hour)
				return &repository.MockAnnouncementRepository{
					GetAnnouncementsFunc: func(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
						return []domain.Announcement{
							{
								ID:             id,
								Title:          "タイトル",
								URL:            "https://example.com/x",
								AvailableFrom:  now,
								AvailableUntil: &until,
							},
						}, nil
					},
				}
			},
			params:   api.AnnouncementsV1ListParams{},
			wantCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response api.AnnouncementsV1List200JSONResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Announcements, 1)
				a := response.Announcements[0]
				assert.Equal(t, "11111111-1111-1111-1111-111111111111", a.Id)
				assert.Equal(t, "タイトル", a.Title)
				assert.Equal(t, "https://example.com/x", a.Url)
				assert.Equal(t, now, a.AvailableFrom)
				assert.NotNil(t, a.AvailableUntil)
				assert.Equal(t, now.Add(48*time.Hour), *a.AvailableUntil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			h := NewHandler(service.NewAnnouncementService(mockRepo))

			request := api.AnnouncementsV1ListRequestObject{Params: tt.params}
			response, err := h.AnnouncementsV1List(context.Background(), request)

			assert.NoError(t, err)
			assert.NotNil(t, response)

			w := httptest.NewRecorder()
			err = response.VisitAnnouncementsV1ListResponse(w)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCode, w.Code)

			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestAnnouncementsV1List_QueryPropagation(t *testing.T) {
	t.Run("クエリパラメータがserviceまで伝搬する", func(t *testing.T) {
		var captured domain.AnnouncementQuery
		sortDesc := api.FoundationV1SortDirection("desc")
		filter := true

		mockRepo := &repository.MockAnnouncementRepository{
			GetAnnouncementsFunc: func(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
				captured = query
				return []domain.Announcement{}, nil
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		request := api.AnnouncementsV1ListRequestObject{
			Params: api.AnnouncementsV1ListParams{
				SortByDate:     &sortDesc,
				FilterIsActive: &filter,
			},
		}
		_, err := h.AnnouncementsV1List(context.Background(), request)
		assert.NoError(t, err)
		assert.Equal(t, domain.SortDirectionDesc, captured.SortByDate)
		assert.True(t, captured.FilterIsActive)
	})

	t.Run("パラメータ未指定時はデフォルト値が使われる", func(t *testing.T) {
		var captured domain.AnnouncementQuery
		mockRepo := &repository.MockAnnouncementRepository{
			GetAnnouncementsFunc: func(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
				captured = query
				return []domain.Announcement{}, nil
			},
		}
		h := NewHandler(service.NewAnnouncementService(mockRepo))

		_, err := h.AnnouncementsV1List(context.Background(), api.AnnouncementsV1ListRequestObject{})
		assert.NoError(t, err)
		assert.Equal(t, domain.SortDirectionAsc, captured.SortByDate)
		assert.False(t, captured.FilterIsActive)
	})
}

func TestAnnouncementsV1List_RepositoryError(t *testing.T) {
	repoErr := errors.New("db unreachable")
	mockRepo := &repository.MockAnnouncementRepository{
		GetAnnouncementsFunc: func(ctx context.Context, query domain.AnnouncementQuery) ([]domain.Announcement, error) {
			return nil, repoErr
		},
	}
	h := NewHandler(service.NewAnnouncementService(mockRepo))

	response, err := h.AnnouncementsV1List(context.Background(), api.AnnouncementsV1ListRequestObject{})
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, response)
}
