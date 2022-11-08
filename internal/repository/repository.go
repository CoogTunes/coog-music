package repository

import (
	"github.com/CoogTunes/coog-music/internal/models"
)

type DatabaseRepo interface {
	//These work
	AddUser(res models.Users) (models.Users, error)
	GetUser(id string) (models.Users, error)
	UpdateUser(user models.Users) (models.Users, error)

	AddArtist(res models.Artist) (models.Artist, error)
	GetArtists() ([]models.Artist, error)
	GetArtistsAndSongs() ([]models.Artist, error)
	UpdateArtist(artist models.Artist) (models.Artist, error)
	GetArtistName(artist_id int) (string, error)

	AddAlbum(album models.Album) (models.Album, error)
	GetAlbums() ([]models.Album, error)

	AddSong(res models.Song) (models.Song, error)
	GetSong(songID string) (models.Song, error)
	GetSongs() ([]models.Song, error)
	UpdateSong(song models.Song) (models.Song, error)

	UpdateSongCount(song models.Song) (models.Song, error)

	GetPlaylists() ([]models.Playlist, error)

	AddSongToAlbum(res models.Song, album models.Album) (models.AlbumSong, error)
	AddSongToPlaylist(song models.Song, playlist models.Playlist) (models.SongPlaylist, error)

	Follow(artistId int, userId int) (models.Followers, error) //add

	Authenticate(email string, password string) error
}
