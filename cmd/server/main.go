package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/PedroMartini98/rss_aggregator_go/config"
	"github.com/PedroMartini98/rss_aggregator_go/internal/api/handler"
	"github.com/PedroMartini98/rss_aggregator_go/internal/api/middleware"
	"github.com/PedroMartini98/rss_aggregator_go/internal/database"
	"github.com/PedroMartini98/rss_aggregator_go/internal/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Error trying to setup database: %v", err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://", "http://"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	router.Mount("/v1", v1Router)

	userHandler := handler.NewUserHandler(dbQueries)
	feedHandler := handler.NewFeedHandler(dbQueries)
	middlewareHandler := middleware.NewMiddlewareHandler(dbQueries)

	//Routes that don't need a handler:
	v1Router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		response.WithJson(w, http.StatusOK, "Server is up")
	})

	// User routes:
	v1Router.Post("/create_user", userHandler.CreateUser)
	v1Router.Get("/user", middlewareHandler.Auth(userHandler.GetUser))

	//Feed routes:
	v1Router.Post("/create_feed", middlewareHandler.Auth(feedHandler.CreateFeed))
	v1Router.Get("/feeds", feedHandler.GetAllFeeds)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + cfg.Port,
	}

	log.Printf("Starting server on port:%v", cfg.Port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
