package repository

import (
	"github.com/CoogTunes/coog-music/internal/models"
)

type DatabaseRepo interface {
	//These work
	AddUser(res models.Users) (models.Users, error)
	GetUser(id string) (models.Users, error)
	AddArtist(res models.Artist) (models.Artist, error)
	GetArtists() ([]models.Artist, error)
	GetArtistsAndSongs() ([]models.Artist, error)
	GetSong(songID string) (models.Song, error)
	AddSong(res models.Song) (models.Song, error)
	GetArtistName(artist_id int) (string, error)
	GetSongs() ([]models.Song, error)
	UpdateSongCount(song models.Song) (models.Song, error)

	AddSongToPlaylist(song models.Song, playlist models.Playlist) (models.SongPlaylist, error)
	AddSongToAlbum(res models.Song, album models.Album) (models.AlbumSong, error)

	GetPlaylists() ([]models.Playlist, error)
	GetAlbums() ([]models.Album, error)

	Authenticate(email string, password string) error
	UpdateUser(user models.Users) (models.Users, error)
	UpdateArtist(artist models.Artist) (models.Artist, error)
	UpdateSong(song models.Song) (models.Song, error)

	Follow(artist models.Artist, user models.Users) (models.Followers, error)
}
