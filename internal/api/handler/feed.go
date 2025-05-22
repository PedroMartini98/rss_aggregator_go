package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PedroMartini98/rss_aggregator_go/internal/database"
	"github.com/PedroMartini98/rss_aggregator_go/internal/response"
	"github.com/google/uuid"
)

type feedHandler struct {
	dbQueries *database.Queries
}

func NewFeedHandler(dbQueries *database.Queries) *feedHandler {
	return &feedHandler{
		dbQueries: dbQueries,
	}
}

func (h *feedHandler) CreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	var p params

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, fmt.Sprintf("unable to parse json: %v", err))
		return
	}

	parsedUrl, err := url.Parse(p.Url)
	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		response.WithError(w, http.StatusBadRequest, fmt.Sprintf("we only accept http urls: %v", err))
		return
	}

	if !strings.Contains(parsedUrl.Host, ".") {
		response.WithError(w, http.StatusBadRequest, fmt.Sprintf("please use a valid url: %v", err))
		return
	}

	feed, err := h.dbQueries.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      p.Name,
		Url:       p.Url,
		UserID:    user.ID,
	})

	if err != nil {
		response.WithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create feed in the database: %v", err))
		return
	}

	response.WithJson(w, http.StatusCreated, feed)
}

func (h *feedHandler) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := h.dbQueries.GetAllFeeds(r.Context())
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get all feeds: %v", err))
		return
	}

	response.WithJson(w, http.StatusOK, feeds)
}

func (h *feedHandler) CreateFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type params struct {
		Feed_id uuid.UUID `json:"feed_id"`
	}

	var p params

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		response.WithError(w, http.StatusInternalServerError, fmt.Sprintf("unable to parse json: %v", err))
		return
	}

	follow, err := h.dbQueries.CreateFollow(r.Context(), database.CreateFollowParams{
		UserID:    user.ID,
		CreatedAt: time.Now(),
		FeedID:    p.Feed_id,
	})

	if err != nil {
		response.WithError(w, http.StatusBadRequest, fmt.Sprintf("unable to create follow: %v", err))
		return
	}
	response.WithJson(w, http.StatusCreated, follow)
}
