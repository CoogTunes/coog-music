package repository

import (
	"github.com/DeLuci/coog-music/internal/models"
)

type DatabaseRepo interface {
	//These work
	AddUser(res models.Users) (models.Users, error)
	GetUser(id string) (models.Users, error)
	AddArtist(res models.Artist) (models.Artist, error)
	GetArtists() ([]models.Artist, error)
	GetSong(songID string) (models.Song, error)
	AddSong(res models.Song) (models.Song, error)
	GetArtistName(artist_id int) (string, error)
	GetSongs() ([]models.Song, error)

	//These' don't work
	AddSongToPlaylist(res models.Song, playlist models.Playlist) error
	AddSongToAlbum(res models.Song, album models.Album) error
}
