package main

// a

import (
	"log"
	"net/http"

	"github.com/PedroMartini98/rss_aggregator_go/config"
	"github.com/PedroMartini98/rss_aggregator_go/internal/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

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

	v1Router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		util.RespondWithJson(w, http.StatusOK, "Server is up")
	})

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
