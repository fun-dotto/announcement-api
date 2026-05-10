package handler

import (
	"testing"
	"time"

	api "github.com/fun-dotto/announcement-api/generated"
	"github.com/fun-dotto/announcement-api/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestToApiAnnouncement(t *testing.T) {
	id := uuid.MustParse("88888888-8888-8888-8888-888888888888")
	from := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	until := from.Add(48 * time.Hour)

	t.Run("全フィールドが変換される", func(t *testing.T) {
		got := toApiAnnouncement(domain.Announcement{
			ID:             id,
			Title:          "タイトル",
			URL:            "https://example.com",
			AvailableFrom:  from,
			AvailableUntil: &until,
		})
		assert.Equal(t, id.String(), got.Id)
		assert.Equal(t, "タイトル", got.Title)
		assert.Equal(t, "https://example.com", got.Url)
		assert.Equal(t, from, got.AvailableFrom)
		assert.NotNil(t, got.AvailableUntil)
		assert.Equal(t, until, *got.AvailableUntil)
	})

	t.Run("AvailableUntilがnilの場合はnilのまま", func(t *testing.T) {
		got := toApiAnnouncement(domain.Announcement{
			ID:            id,
			Title:         "x",
			URL:           "https://example.com",
			AvailableFrom: from,
		})
		assert.Nil(t, got.AvailableUntil)
	})
}

func TestToDomainAnnouncementQuery(t *testing.T) {
	t.Run("両パラメータ未指定時はデフォルト", func(t *testing.T) {
		got := toDomainAnnouncementQuery(api.AnnouncementsV1ListParams{})
		assert.Equal(t, domain.SortDirectionAsc, got.SortByDate)
		assert.False(t, got.FilterIsActive)
	})

	t.Run("SortByDate=descが反映される", func(t *testing.T) {
		desc := api.FoundationV1SortDirection("desc")
		got := toDomainAnnouncementQuery(api.AnnouncementsV1ListParams{SortByDate: &desc})
		assert.Equal(t, domain.SortDirectionDesc, got.SortByDate)
	})

	t.Run("FilterIsActive=trueが反映される", func(t *testing.T) {
		v := true
		got := toDomainAnnouncementQuery(api.AnnouncementsV1ListParams{FilterIsActive: &v})
		assert.True(t, got.FilterIsActive)
	})

	t.Run("FilterIsActive=falseがそのまま反映される", func(t *testing.T) {
		v := false
		got := toDomainAnnouncementQuery(api.AnnouncementsV1ListParams{FilterIsActive: &v})
		assert.False(t, got.FilterIsActive)
	})
}

func TestToDomainAnnouncementFromRequest(t *testing.T) {
	id := uuid.MustParse("99999999-9999-9999-9999-999999999999")
	from := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	until := from.Add(24 * time.Hour)

	t.Run("全フィールドがドメインに反映される", func(t *testing.T) {
		got := toDomainAnnouncementFromRequest(id, api.AnnouncementRequest{
			Title:          "t",
			Url:            "https://example.com",
			AvailableFrom:  from,
			AvailableUntil: &until,
		})
		assert.Equal(t, id, got.ID)
		assert.Equal(t, "t", got.Title)
		assert.Equal(t, "https://example.com", got.URL)
		assert.Equal(t, from, got.AvailableFrom)
		assert.Equal(t, &until, got.AvailableUntil)
	})

	t.Run("AvailableUntilがnilの場合はnilのまま", func(t *testing.T) {
		got := toDomainAnnouncementFromRequest(id, api.AnnouncementRequest{
			Title:         "t",
			Url:           "https://example.com",
			AvailableFrom: from,
		})
		assert.Nil(t, got.AvailableUntil)
	})
}
