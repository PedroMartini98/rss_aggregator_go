package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PedroMartini98/rss_aggregator_go/internal/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("Falied to load PORT from enviroment")
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
		Addr:    ":" + portString,
	}

	log.Printf("Starting server on port:%v", portString)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
