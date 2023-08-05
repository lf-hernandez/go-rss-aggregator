package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/lf-hernandez/go-rss-aggregator/graph"
	"github.com/lf-hernandez/go-rss-aggregator/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	Database *database.Queries
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
		Database: db,
	}

	go startScraping(db, 10, time.Minute)

	port := os.Getenv("PORT")
	if port == "" {
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

	v1Router.Get("/posts", apiConfiguration.middlewareAuth(apiConfiguration.handlerGetPostsByUser))

	router.Mount("/v1", v1Router)

	v2Router := chi.NewRouter()
	gqlHandler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Database: db,
	}}))

	v2Router.Handle("/", playground.Handler("GraphQL playground", "/v2/query"))
	v2Router.Handle("/query", gqlHandler)

	router.Mount("/v2", v2Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server starting on port %v", port)

	serverError := server.ListenAndServe()
	if serverError != nil {
		log.Fatal(serverError)
	}
}
