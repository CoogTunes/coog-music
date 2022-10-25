package repository

import "github.com/DeLuci/coog-music/internal/models"

type DatabaseRepo interface {
	AddSong(res models.Song)
	AddUser(res models.Users)
	AddSongToPlaylist(res models.Song, playlist models.Playlist)
	PlaySong(res models.Song)
	AddSongToAlbum(res models.Song, album models.Album)
}
