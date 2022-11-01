package dbrepo

import (
	"database/sql"
	"log"

	"github.com/DeLuci/coog-music/internal/models"
)

// USERS
func (m *postgresDBRepo) AddUser(res models.Users) (models.Users, error) {
	var user models.Users

	query := "insert into Users (username, password, first_name, last_name, gender, admin) values ($1, $2, $3, $4, $5, $6) RETURNING *"

	row := m.DB.QueryRow(query, res.Username, res.Password, res.First_name, res.Last_name, res.Gender, res.Admin)

	err := row.Scan(&user.User_id, &user.Username, &user.First_name, &user.Last_name, &user.Gender, &user.Password, &user.Admin)
	if err != nil {
		log.Println(err)
	}
	return user, nil
}

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

// ARTISTS
func (m *postgresDBRepo) AddArtist(res models.Artist) (models.Artist, error) {
	var artist models.Artist

	query := "insert into Artist (name, artist_id, location, join_date) values ($1, $2, $3, to_date($4, 'YYYY-MM-DD')) RETURNING *"
	row := m.DB.QueryRow(query, res.Name, res.Artist_id, res.Location, res.Join_date)

	err := row.Scan(&artist.Name, &artist.Artist_id, &artist.Location, &artist.Join_date)
	if err != nil {
		log.Println(err)
	}
	return artist, nil
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
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var artist models.Artist

		rows.Scan(&artist.Name, &artist.Artist_id, &artist.Location, &artist.Join_date)

		if err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}
	return artists, nil
}

func (m *postgresDBRepo) GetArtistName(artist_id int) (str string, err error) {
	var song models.Song
	query := `SELECT name FROM Artist as A, Song as S WHERE A.artist_id = S.artist_id AND A.artist_id = $1`
	row := m.DB.QueryRow(query, artist_id)
	err2 := row.Scan(&song.Artist_name)
	if err2 != nil {
		log.Println(err)
	}

	return song.Artist_name, err
}

// SONGS
func (m *postgresDBRepo) AddSong(res models.Song) (models.Song, error) {

	var song models.Song

	query := `insert into song (title, artist_id, release_date, duration, album_id, total_plays)
				 select $1, ar.artist_id, to_date($2, 'YYYY-MM-DD'), $3, al.album_id, 0 from artist as ar, album as al where ar.artist_id = $4 AND al.album_id = $5 RETURNING *`
	row := m.DB.QueryRow(query, res.Title, res.Release_date, res.Duration, res.Artist_id, res.Album_id)

	err := row.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album_id, &song.Total_plays)
	if err != nil {
		log.Println(err)
	}
	return song, nil
}

func (m *postgresDBRepo) GetSongs() ([]models.Song, error) {
	var songs []models.Song
	// probably need to add a where statement and get rid of *
	query := "SELECT * FROM song"

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
		var song models.Song

		rows.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album_id, &song.Total_plays)

		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func (m *postgresDBRepo) GetSong(songID string) (models.Song, error) {

	var song models.Song

	query := "select * from song where song_id = $1"

	row := m.DB.QueryRow(query, songID)
	log.Println("row", row)
	log.Println(row.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album, &song.Total_plays))

	return song, nil

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

func (m *postgresDBRepo) AddSongToAlbum(res models.Song, album models.Album) (models.Album, error) {
	// query := "select song from song where title == $1"
	// add_query := "insert into song(album) values ($1)"

	var albums models.Album

	// query := `insert into album(name, artist_id, date_added, song_id)
	// select $1, $2, $3, to_date($4, 'YYY-MM-DD'), song_id from song where song_id = $5 returning *`

	// row := m.DB.QueryRow(query, res.Artist_id, res.Date_added, res.Song_id)

	// err := row.Scan(&res.Name, &res.Artist_id, &res.Album_id, &res.Date_added, song_id)
	// if err != nil {
	// 	log.Println(err)
	// }

	return albums, nil
}

func (m *postgresDBRepo) GetPlaylists() ([]models.Playlist, error) {
	var playlists []models.Playlist

	query := "SELECT * FROM PLAYLIST"

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var playlist models.Playlist

		rows.Scan(&playlist.User_id, &playlist.Name, &playlist.Playlist_id)

		if err != nil {
			return nil, err
		}

		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (m *postgresDBRepo) GetAlbums() ([]models.Album, error) {
	var albums []models.Album

	query := "SELECT * FROM ALBUM"

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var album models.Album

		rows.Scan(
			&album.Name,
			&album.Album_id,
			&album.Artist_id,
			&album.Date_added)

		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}
	return albums, nil
}

//inset, select, update, delete
//album, playlist
