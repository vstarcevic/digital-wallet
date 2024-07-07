package api

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Routes(cfg *Config) http.Handler {
	router := chi.NewRouter()

	// specify who is allowed to connect
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Use(middleware.Heartbeat("/health"))

	router.Get("/time", cfg.getTime)
	router.Post("/user", cfg.createUser)

	router.Route("/", func(r chi.Router) {
		r.Handle("/", http.RedirectHandler("/new/", http.StatusPermanentRedirect))
	})

	return router
}
