package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"backend/config"
	"backend/controllers"
	"backend/models"
)

func setupTestEcho() *echo.Echo {
	e := echo.New()

	// Optional: insert dummy data here if needed
	config.DB.Create(&models.MusicPreference{
		UserID:    1,
		TrackID:   1,
		IsLiked:   false,
		PlayCount: 0,
	})

	return e
}

func TestLikeTrack(t *testing.T) {
	e := setupTestEcho()
	req := httptest.NewRequest(http.MethodPost, "/like/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/like/:track_id")
	c.SetParamNames("track_id")
	c.SetParamValues("1")
	c.Set("uid", uint(1))

	if assert.NoError(t, controllers.LikeTrack(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"is_liked":true`)
	}
}

func TestDislikeTrack(t *testing.T) {
	e := setupTestEcho()
	req := httptest.NewRequest(http.MethodPost, "/dislike/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/dislike/:track_id")
	c.SetParamNames("track_id")
	c.SetParamValues("1")
	c.Set("uid", uint(1))

	if assert.NoError(t, controllers.DislikeTrack(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"is_liked":false`)
	}
}

func TestGetLikedTracks(t *testing.T) {
	e := setupTestEcho()
	req := httptest.NewRequest(http.MethodGet, "/liked", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("uid", uint(1))

	if assert.NoError(t, controllers.GetLikedTracks(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// If a track is liked, we should see it here
	}
}

func TestGetTrackPlayCount(t *testing.T) {
	e := setupTestEcho()
	req := httptest.NewRequest(http.MethodGet, "/playcount/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/playcount/:track_id")
	c.SetParamNames("track_id")
	c.SetParamValues("1")
	c.Set("uid", uint(1))

	if assert.NoError(t, controllers.GetTrackPlayCount(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `0`) // Assuming starting value
	}
}

func TestIncrementTrackPlayCount(t *testing.T) {
	e := setupTestEcho()
	req := httptest.NewRequest(http.MethodPost, "/playcount/1/increment", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/playcount/:track_id/increment")
	c.SetParamNames("track_id")
	c.SetParamValues("1")
	c.Set("uid", uint(1))

	if assert.NoError(t, controllers.IncrementTrackPlayCount(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"play_count":1`)
	}
}
