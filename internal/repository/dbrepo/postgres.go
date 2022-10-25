package dbrepo

import "github.com/DeLuci/coog-music/internal/models"

func (n *postgresDBRepo) AddUser(res models.Users) error {
	query := "insert into users (username, password) values ($1, $2)"

	_, err := m.DB.Exec(query, res.Username, res.Password)
	if err != nil {
		return err
	}

	return nil
}

func (n *postgresDBRepo) AddSong(res models.Song) error {
	query := "insert into song (title, artist_name) values ($1, $2)"

	_, err := m.Db.Exec(query, res.Artist_name, res.Title)
	if err != nil {
		return err
	}

	return nil
}

//TODO: ADD LINKING TABLES AND USE THEM TO GRAB THE OTHER STUFF

func (n *postgresDBRepo) AddSongToPlaylist(song models.Song, playlist models.Playlist) error {
	query := "insert into playlist (playlist.playlist_id, playlist.songs) values($1, $2)"

	_, err := m.Db.Exec(query, playlist.Playlist_id, song.N)
	if err != nil {
		return err
	}

	return nil
}


func (n *postgresDBRepo) PlaySong(res models.Song) error {
	query := "select song_id from song where title == $1"

	_, err := m.DB.Exec(query, res.song_id)
	if err != nil {
		return err
	}
	//generate a session id for songplay

	return nil
}

func (n* postgresDBRepo) AddSongToAlbum(res models.Song, album models.Album) error{
	query := "select song from song where title == $1"
	add_query := "insert into song(album) values ($2)"

	_, err := m.Db.Exec(query, res.title)
	if err != nil{
		return err
	}
	_, err := m.Db.Exec(add_query, album.name)
	if err != nil{
		return err
	}

	return nil
}

//play song (select and songplay session)
//add song to album (artist thing
