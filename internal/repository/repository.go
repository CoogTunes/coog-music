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
	AddPlaylist(res models.Playlist) (models.Playlist, error)
	AddPlaylistSong(num int, playlistID int) error
	AddSongForAlbum(res models.Song) error
	GetTopSongs() ([]models.Song, error)
	GetUser(id string) (models.Users, error)
	//GetSong(songID string) (models.Song, error)
	//GetSongs() ([]models.Song, error)
	GetPlaylists(userId int) ([]models.Playlist, error)
	GetAlbums() ([]models.Album, error)
	//GetArtistName(artist_id int) (string, error)
	GetArtists() ([]models.Artist, error)
	//GetArtistsAndSongs() ([]models.Artist, error)
	//SearchPlaylists(playlist_name string) ([]models.Playlist, error)
	//SearchAlbums(album_name string) ([]models.Album, error)
	//SearchArtists(artist_name string) ([]models.Artist, error)
	GetSongsByName(song_name string) ([]models.Song, error)
	GetSongsFromPlaylist(playlist_id int) ([]models.Song, error)
	GetSongsFromArtist(artist_name string) (map[string][]models.Song, error)
	GetSongsFromArtistByID(artist_name string, artistID int) (map[string][]models.Song, error)
	GetSongsFromAlbum(album_name string) ([]models.Song, error)
	GetSongsFromAlbumByID(album_name string, albumID int) ([]models.Song, error)
	GetSongsForLikePage(userId int) ([]models.Song, error)
	GetInitialUsersReport() ([]models.UserReport, error)
	GetNumberOfUsers() (models.Users, error)
	GetNumberOfSongs() (models.Song, error)
	GetNumberOfArtists() (models.Users, error)
	GetNumberOfPlaylists() (models.Playlist, error)

	Follow(artistId int, userId int) (models.Followers, error) //add
	Authenticate(email string, password string) (models.Users, error)
	AddOrUpdateLikeValue(islike bool, songId int, userId int) error
	SendUpdatedLikeValue(songID int) (models.Song, error)
	UpdateMessages(user_id int) ([]models.Messages, error)
	// UpdateMessage(user_id int)(error)

	UpdateUser(user models.Users) (models.Users, error)
	// UpdateSong(song models.Song) (models.Song, error)
	UpdateSongCount(id int) (models.Song, error)
	UpdateArtist(artist models.Artist) (models.Artist, error)

	RemoveUser(user_id int) error
	RemoveSong(song_id int) error
	RemoveAlbum(album_id int) error
	RemovePlaylist(playlist_id int) error
	RemoveArtist(artist_id int) error
	RemoveSongFromAlbum(song_id int, album_id int) error
	RemoveSongFromPlaylist(song_id int, playlist_id int) error

	GetLikesReport(minLikes int, maxLikes int) ([]models.Song, error)
	GetUsersReport(minDate string, maxDate string) ([]models.UserReport, error)
	GetArtistReport(minDate string, maxDate string) ([]models.ArtistReport, error)
	GetPlaysReport(minDate string, maxDate string, min_plays int, max_plays int) ([]models.Song, error)
}
