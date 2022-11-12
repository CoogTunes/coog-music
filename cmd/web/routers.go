package main

import (
	"net/http"

	// "github.com/DeLuci/coog-music/internal/config"
	"github.com/CoogTunes/coog-music/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	// Commenting out to get POST requests to work.
	// mux.Use(NoSurf)
	// mux.Use(SessionLoad)

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.AllowContentType("application/json", "application/x-www-form-urlencoded", "multipart/form-data"))
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))
	// These work
	mux.Post("/user", handlers.Repo.AddUser)
	mux.Get("/user/{id}", handlers.Repo.GetUser)
	mux.Put("/user", handlers.Repo.UpdateUser)

	//mux.Post("/artist", handlers.Repo.AddArtist)
	mux.Get("/artists", handlers.Repo.GetArtists)
	mux.Put("/artist", handlers.Repo.UpdateArtist)
	//mux.Get("/artists-songs", handlers.Repo.GetArtistsAndSongs)

	//mux.Post("/album", handlers.Repo.AddAlbum)
	mux.Get("/albums", handlers.Repo.GetAlbums)

	//mux.Post("/song", handlers.Repo.AddSong)
	//mux.Get("/song/{id}", handlers.Repo.GetSong)
	//mux.Get("/songs", handlers.Repo.GetSongs)
	//mux.Put("/song", handlers.Repo.UpdateSong)

	//mux.Put("/songCount/{id}", handlers.Repo.UpdateSongCount)

	mux.Get("/playlists", handlers.Repo.GetPlaylists)

	mux.Post("/album/{songid}", handlers.Repo.AddSongToAlbum)
	mux.Post("/song/{playlistId}", handlers.Repo.AddSongToPlaylist)

	mux.Post("/follow", handlers.Repo.Follow)

	// Register and Login Handlers
	mux.Get("/login", handlers.Repo.GetLogin)
	mux.Post("/login", handlers.Repo.PostLogin)

	// Logout Handler
	mux.Get("/logout", handlers.Repo.LogOut)

	// Home page handlers
	mux.Get("/", handlers.Repo.GetHome)

	// Upload handlers
	mux.Post("/upload", handlers.Repo.UploadFile)

	// Playlist
	mux.Get("/playlist/Search/{playlistSearchInfo}", handlers.Repo.PlaylistSearch)
	mux.Post("/addPlaylist", handlers.Repo.InsertPlaylist)
	// Profile page handlers
	mux.Get("/profile", handlers.Repo.GetProfile)
	// Need to finish handlers and maybe adjust routing

	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) // helps use css/js

	return mux
}
