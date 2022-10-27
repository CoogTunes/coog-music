package repository

import (
	"github.com/DeLuci/coog-music/internal/models"
)

type DatabaseRepo interface {
	GetArtists() ([]models.Artist, error)
	GetUsers() ([]models.Users, error)
	AddSong(res models.Song) error
	AddUser(res models.Users) (models.Users, error)
	AddSongToPlaylist(res models.Song, playlist models.Playlist) error
	GetSong(songID string) (models.Song, error)
	AddSongToAlbum(res models.Song, album models.Album) error
}
