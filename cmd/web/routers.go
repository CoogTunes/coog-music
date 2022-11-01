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

	// This one works
	mux.Get("/artists", handlers.Repo.GetArtists)
	mux.Get("/users", handlers.Repo.GetUsers)
	mux.Get("/login", handlers.Repo.GetLogin)
	mux.Post("/login", handlers.Repo.PostLogin)
	// Need to finish handlers and maybe adjust routing
	mux.Post("/song", handlers.Repo.AddSong)
	mux.Post("/user", handlers.Repo.AddUser)
	mux.Post("/song/{playlistId}", handlers.Repo.AddSongToPlaylist)
	mux.Get("/song/{id}", handlers.Repo.GetSong)
	mux.Post("/album/{songid}", handlers.Repo.AddSongToAlbum)

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) // helps use css/js
	
	return mux
}
