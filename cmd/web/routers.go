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

	// Commenting out to get POST requests to work.
	// mux.Use(NoSurf)
	// mux.Use(SessionLoad)

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.AllowContentType("application/json", "application/x-www-form-urlencoded", "multipart/form-data"))

	// These work
	mux.Get("/users", handlers.Repo.GetUsers)
	mux.Get("/song/{id}", handlers.Repo.GetSong)

	//might work
	mux.Get("/artists", handlers.Repo.GetArtists)

	// Need to finish handlers and maybe adjust routing
	mux.Post("/song", handlers.Repo.AddSong)
	mux.Post("/user", handlers.Repo.AddUser)
	mux.Post("/song/{playlistId}", handlers.Repo.AddSongToPlaylist)
	mux.Post("/album/{songid}", handlers.Repo.AddSongToAlbum)

	return mux
}
