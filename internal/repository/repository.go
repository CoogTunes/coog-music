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
	AddSong(res models.Song, artist_id int) (models.Song, error)
	GetArtistName(artist_id int) (string, error)

	//These' don't work
	AddSongToPlaylist(playlist models.Playlist, users models.Users, songs models.Song) (models.Playlist, error)
	AddSongToAlbum(res models.Album, song_id int) (models.Album, error)

	GetPlaylists()([]models.Playlist, error)
	GetAlbums()([]models.Album, error)
}
