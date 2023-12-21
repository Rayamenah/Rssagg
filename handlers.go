package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Rayamenah/rssagg/internal/auth"
	"github.com/Rayamenah/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type apiConfig struct {
	DB *database.Queries
}

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		writeErr(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		writeErr(w, 400, fmt.Sprintf("couldnt create user: %v", err))
		return
	}

	writeJSON(w, 201, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		writeErr(w, 403, fmt.Sprintf("Auth error: %v", err))
		return
	}
	user, err = apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		writeErr(w, 400, fmt.Sprintf("couldn't get user: %v", err))
		return
	}
	writeJSON(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		writeErr(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})
	if err != nil {
		writeErr(w, 400, fmt.Sprintf("couldn't create user: %v", err))
		return
	}

	writeJSON(w, 201, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		writeErr(w, 400, fmt.Sprintf("Error parsing JSON %v", err))
		return
	}
	feed, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		writeErr(w, 400, fmt.Sprintf("couldnt create feed follow: %v", err))
		return
	}

	writeJSON(w, 201, feed)
}
func (apiCfg *apiConfig) handleGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		writeErr(w, 400, fmt.Sprintf("couldnt create feed follow: %v", err))
		return
	}

	writeJSON(w, 201, databaseFeedFollowsToFollows(feedFollows))
}

// func (apiCfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {

// 	feeds, err := apiCfg.DB.GetFeeds(r.Context())
// 	if err != nil {
// 		writeErr(w, 400, fmt.Sprintf("couldnt get feeds: %v", err))
// 		return
// 	}

// 	writeJSON(w, 201, databaseFeedsToFeed(feeds))
// }

func (apiCfg *apiConfig) handleDeleteFeedfollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	FeedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		writeErr(w, 400, fmt.Sprintf("couldnt parse feed follow id: %v", err))
		return
	}
	userIDstr := chi.URLParam(r, "user_id")
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		writeErr(w, 400, fmt.Sprintf("couldnt parse feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.UnfollowFeed(r.Context(), database.UnfollowFeedParams{
		ID:     FeedFollowID,
		UserID: userID,
	})
	if err != nil {
		writeErr(w, 400, fmt.Sprintf("couldnt delete feed follow: %v", err))
		return
	}
	writeJSON(w, 200, struct{}{})
}

func (apiCfg *apiConfig) handleGetPostForUser(w http.ResponseWriter, r *http.Request, user database.User) {

	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		writeErr(w, 400, fmt.Sprintf("couldntget posts: %v", err))
		return
	}
	writeJSON(w, 200, databasePostsToPosts(posts))
}
