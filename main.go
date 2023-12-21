package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Rayamenah/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func writeErr(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error:", msg)
	}
	type errResp struct {
		Error string `json:"error"`
	}

	writeJSON(w, code, errResp{Error: msg})
}

func writeJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal Json response %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found in the env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("PORT is not found in the env")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cant conect to database", err)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	go StartScrapping(apiCfg.DB, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	//routes
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleErr)
	v1Router.Get("/err", handleErr)
	v1Router.Post("/users", apiCfg.handleCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handleGetUser))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	// v1Router.Get("/feeds", apiCfg.handleGetFeeds)

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handleGetPostForUser))

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handleGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handleDeleteFeedfollow))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Printf("server starting on port %v", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func handleErr(w http.ResponseWriter, r *http.Request) {
	writeErr(w, 400, "something went wrong")
}
