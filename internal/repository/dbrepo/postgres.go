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
			log.Println(err)
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

	query := "insert into Users (username, password, first_name, last_name, gender, admin) values ($1, $2, $3, $4, $5, $6) RETURNING *"

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

func (m *postgresDBRepo) AddSong(res models.Song, artist_id int) (models.Song, error) {

	var song models.Song

	query := `insert into song (title, artist_id, release_date, duration, album, total_plays)
				 select $1, artist_id, to_date($2, 'YYYY-MM-DD'), $3, $4, 0 from artist where artist_id = $5 RETURNING *`
	row := m.DB.QueryRow(query, res.Title, res.Release_date, res.Duration, res.Album, artist_id)

	err := row.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album, &song.Total_plays)
	if err != nil {
		log.Println(err)
	}
	return song, nil
}

func (m *postgresDBRepo) GetArtistName(artist_id int) (str string, err error) {
	var song models.Song
	query := `SELECT name FROM Artist as A, Song as S WHERE A.artist_id = S.artist_id AND A.artist_id = $1`
	row := m.DB.QueryRow(query, artist_id)
	err2 := row.Scan(&song.Artist_id)
	if err2 != nil {
		log.Println(err)
	}

	return song.Artist_id, err
}

func (m *postgresDBRepo) AddSongToPlaylist(playlist models.Playlist, users models.Users, songs models.Song) (models.Playlist, error) {
	query := `insert into playlist (user_id, name, playlist_id, playlist_length, song_id)
				select users.user_id, song.title, $1, $2, song.song_id from users, song where users.user id = $3 and song.song_id = $4`
	var pl models.Playlist

	row := m.DB.QueryRow(query, playlist.Playlist_id, playlist.Playlist_length, users.User_id, songs.Song_id)
	err := row.Scan(&playlist.User_id, &playlist.Name, &playlist.Playlist_id, &playlist.Playlist_length, &playlist.Song_id)

	if err != nil {
		log.Println(err)
	}

	return pl, nil
}

func (m *postgresDBRepo) GetSong(songID string) (models.Song, error) {

	var song models.Song

	query := "select * from song where song_id = $1"

	row := m.DB.QueryRow(query, songID)
	log.Println("row", row)
	log.Println(row.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album, &song.Total_plays))

	return song, nil

}

func (m *postgresDBRepo) AddSongToAlbum(res models.Album, song_id int) (models.Album, error) {

	var album models.Album

	query := `insert into album(name, artist_id, date_added, song_id)
	select $1, $2, $3, to_date($4, 'YYY-MM-DD'), song_id from song where song_id = $5 returning *`

	row := m.DB.QueryRow(query, res.Name, res.Artist_id, res.Date_added, res.Song_id)

	err := row.Scan(&res.Name, &res.Artist_id, &res.Album_id, &res.Date_added, song_id)
	if err != nil {
		log.Println(err)
	}

	return album, nil
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

		rows.Scan(&playlist.User_id, &playlist.Name, &playlist.Playlist_id, &playlist.Playlist_length, &playlist.Song_id)

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
			&album.Date_added,
			&album.Song_id)

		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}
	return albums, nil
}

//inset, select, update, delete
//album, playlist
