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
	mux.Post("/user", handlers.Repo.AddUser)
	mux.Get("/user/{id}", handlers.Repo.GetUser)
	mux.Put("/user", handlers.Repo.UpdateUser)

	mux.Post("/artist", handlers.Repo.AddArtist)
	mux.Get("/artists", handlers.Repo.GetArtists)

	mux.Get("/song/{id}", handlers.Repo.GetSong)
	mux.Post("/song", handlers.Repo.AddSong)
	mux.Get("/songs", handlers.Repo.GetSongs)

	mux.Put("/songCount/{id}", handlers.Repo.UpdateSongCount)

	mux.Get("/playlists", handlers.Repo.GetPlaylists)

	mux.Get("/albums", handlers.Repo.GetAlbums)

	mux.Get("/login", handlers.Repo.GetLogin)
	mux.Post("/login", handlers.Repo.PostLogin)
	// Need to finish handlers and maybe adjust routing
	mux.Post("/song/{playlistId}", handlers.Repo.AddSongToPlaylist)
	mux.Post("/album/{songid}", handlers.Repo.AddSongToAlbum)

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) // helps use css/js

	return mux
}
