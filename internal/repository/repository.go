package repository

import (
	"github.com/CoogTunes/coog-music/internal/models"
)

type DatabaseRepo interface {
	//These work
	AddUser(res models.Users) (models.Users, error)
	AddArtist(res models.Artist) error
	AddAlbum(album models.Album) (models.Album, error)
	AddSongToAlbum(res models.Song, album models.Album) (models.AlbumSong, error)
	AddSongToPlaylist(song models.Song, playlist models.Playlist) (models.SongPlaylist, error)
	AddSong(res models.Song) error
	AddPlaylist(res models.Playlist) error
	AddSongForAlbum(res models.Song) error
	GetUser(id string) (models.Users, error)
	//GetSong(songID string) (models.Song, error)
	//GetSongs() ([]models.Song, error)
	GetPlaylists() ([]models.Playlist, error)
	GetAlbums() ([]models.Album, error)
	//GetArtistName(artist_id int) (string, error)
	GetArtists() ([]models.Artist, error)
	//GetArtistsAndSongs() ([]models.Artist, error)
	//SearchPlaylists(playlist_name string) ([]models.Playlist, error)
	//SearchAlbums(album_name string) ([]models.Album, error)
	//SearchArtists(artist_name string) ([]models.Artist, error)
	GetSongsByName(song_name string) ([]models.Song, error)
	GetSongsFromPlaylist(playlist_name string) ([]models.Song, error)
	GetSongsFromArtist(artist_name string) ([]models.Song, error)
	GetSongsFromAlbum(album_name string) ([]models.Song, error)

	GetNumberOfUsers() (models.Users, error)
	GetNumberOfSongs() (models.Song, error)
	GetNumberOfArtists() (models.Users, error)
	GetNumberOfPlaylists() (models.Playlist, error)

	Follow(artistId int, userId int) (models.Followers, error) //add
	Authenticate(email string, password string) (models.Users, error)
	AddLikeToSong(song_id int) error

	UpdateUser(user models.Users) (models.Users, error)
	UpdateSong(song models.Song) (models.Song, error)
	UpdateSongCount(song models.Song) (models.Song, error)
	UpdateArtist(artist models.Artist) (models.Artist, error)

	RemoveUser(user_id int) error
	RemoveSong(song_id int) error
	RemoveAlbum(album_id int) error
	RemovePlaylist(playlist_id int) error
	RemoveArtist(artist_id int) error
	RemoveSongFromAlbum(song_id int, album_id int) error
	RemoveSongFromPlaylist(song_id int, playlist_id int) error

	GetLikesReport(minLikes int, maxLikes int, minDislikes int, maxDislikes int) ([]models.LikesReport, error)
}
