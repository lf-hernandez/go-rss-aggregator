package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/lf-hernandez/go-rss-aggregator/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	dotEnvError := godotenv.Load()
	if dotEnvError != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not found")
	}

	sqlConnection, sqlError := sql.Open("postgres", dbURL)
	if sqlError != nil {
		log.Fatal("Error connecting to database: ", sqlError)
	}

	db := database.New(sqlConnection)
	apiConfiguration := apiConfig{
		DB: db,
	}

	go startScraping(db, 10, time.Minute)

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not found")
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)

	v1Router.Get("/error", handlerError)

	v1Router.Post("/users", apiConfiguration.handlerCreateUser)
	v1Router.Get("/users", apiConfiguration.middlewareAuth(apiConfiguration.handlerGetUser))

	v1Router.Post("/feeds", apiConfiguration.middlewareAuth(apiConfiguration.handlerCreateFeed))
	v1Router.Get("/feeds", apiConfiguration.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiConfiguration.middlewareAuth(apiConfiguration.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiConfiguration.middlewareAuth(apiConfiguration.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiConfiguration.middlewareAuth(apiConfiguration.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	serverError := server.ListenAndServe()
	if serverError != nil {
		log.Fatal(serverError)
	}
}
