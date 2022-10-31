package dbrepo

import (
	"database/sql"
	"log"

	"github.com/DeLuci/coog-music/internal/models"
)

// For logging in?
func (m *postgresDBRepo) GetUser(User_id string) (models.Users, error) {

	var user models.Users

	query := "SELECT * FROM Users WHERE user_id = $1"
	rows := m.DB.QueryRow(query, User_id)

	err := rows.Scan(&user.User_id, &user.Username, &user.First_name, &user.Last_name, &user.Gender, &user.Password, &user.Admin)
	if err != nil {

	}

	return user, nil
}

// For searching artists?
func (m *postgresDBRepo) GetArtists() ([]models.Artist, error) {
	var artists []models.Artist
	// probably need to add a where statement and get rid of *
	query := "SELECT * FROM artist"

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var artist models.Artist

		rows.Scan(&artist.Name, &artist.Artist_id, &artist.Location, &artist.Join_date, &artist.Admin)

		if err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}
	return artists, nil
}

func (m *postgresDBRepo) AddUser(res models.Users) (models.Users, error) {
	var user models.Users

	query := "insert into Users (username, password, first_name, last_name, gender, admin) values ($1, $2, $3, $4, $5, $6)"

	row := m.DB.QueryRow(query, res.Username, res.Password, res.First_name, res.Last_name, res.Gender, res.Admin)

	err := row.Scan(&user.User_id, &user.Username, &user.First_name, &user.Last_name, &user.Gender, &user.Password, &user.Admin)
	if err != nil {
		log.Println(err)
	}
	return user, nil
}

func (m *postgresDBRepo) AddArtist(res models.Artist) (models.Artist, error) {
	var artist models.Artist

	query := "insert into Artist (name, artist_id, location, join_date, admin) values ($1, $2, $3, to_date($4, 'YYYY-MM-DD'), $5) RETURNING *"
	row := m.DB.QueryRow(query, res.Name, res.Artist_id, res.Location, res.Join_date, res.Admin)

	err := row.Scan(&artist.Name, &artist.Artist_id, &artist.Location, &artist.Join_date, &artist.Admin)
	if err != nil {
		log.Println(err)
	}
	return artist, nil
}

func (m *postgresDBRepo) AddSong(res models.Song) error {
	query := "insert into song (title, artist_name) values ($1, $2)"

	_, err := m.DB.Exec(query, res.Artist_name, res.Title)
	if err != nil {
		return err
	}

	return nil
}

//TODO: ADD LINKING TABLES AND USE THEM TO GRAB THE OTHER STUFF

func (m *postgresDBRepo) AddSongToPlaylist(song models.Song, playlist models.Playlist) error {
	query := "insert into playlist (playlist.playlist_id, playlist.songs) values($1, $2)"

	_, err := m.DB.Exec(query, playlist.Playlist_id, song.Song_id)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) GetSong(songID string) (models.Song, error) {

	var song models.Song

	query := "select * from song where song_id = $1"

	row := m.DB.QueryRow(query, songID)
	log.Println("row", row)
	log.Println(row.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Artist_name, &song.Album, &song.Total_plays))

	return song, nil

}

func (m *postgresDBRepo) AddSongToAlbum(res models.Song, album models.Album) error {
	query := "select song from song where title == $1"
	add_query := "insert into song(album) values ($1)"

	_, err := m.DB.Exec(query, res.Title)
	if err != nil {
		return err
	}
	_, err2 := m.DB.Exec(add_query, album.Name)
	if err2 != nil {
		return err2
	}

	return nil
}

//play song (select and songplay session)
//add song to album (artist thing
