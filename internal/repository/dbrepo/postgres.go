package dbrepo

import "github.com/DeLuci/coog-music/internal/models"

func (n *postgresDBRepo) AddUser(res models.Users) err {
	query := "insert into users (username, password) values ($1, $2)"

	_, err := m.DB.Exec(query, res.Username, res.Password)
	if err != nil {
		return err
	}

	return nil
}

func (n *postgresDBRepo) AddSong(res models.Song) err{
	query := "insert into song (title, artist_name) values ($1, $2)"

	_, err := m.Db.Exec(query, res.Artist_name, res.Title)
	if err != nil{
		return err
	}

	return nil
}

func (n* postgresDBRepo) AddSongToPlaylist(res models.Song, models.Playlist) err{
	query := "insert into playlist (playlist.songs) values($1)"

	_, err := m.Db. Exec(query, res.song_id)
	if err != nil{
		return err
	}

	return nil
}

//play song (select and songplay session)
//add song to album (artist thing
