package main

import (
	"net/http"

	// "github.com/DeLuci/coog-music/internal/config"
	"github.com/DeLuci/coog-music/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.GetArtists)

	return mux
}
