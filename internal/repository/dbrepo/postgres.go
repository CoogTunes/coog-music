package dbrepo

import (
	"context"
	"database/sql"
	"github.com/DeLuci/coog-music/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
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
		log.Println(err)
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

func (m *postgresDBRepo) Authenticate(email string, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var hashedPwd string
	row := m.DB.QueryRowContext(ctx, "select password from users where username = $1", email)
	err := row.Scan(&hashedPwd)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return err
	} else if err != nil {
		return err
	}

	return nil
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
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var song models.Song

		rows.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album_id, &song.Total_plays)

		if err != nil {
			log.Println(err)
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

	//maybe call update playcount
	return song, nil

}

func (m *postgresDBRepo) UpdateSongCount(song models.Song) (models.Song, error){
	var songs models.Song

	query := "UPDATE SONG SET total_plays = total_plays + 1"
	
	row := m.DB.QueryRow(query)
	log.Println("row", row)
	log.Println(row.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album, &song.Total_plays))

	return songs, nil

}

//TODO: ADD LINKING TABLES AND USE THEM TO GRAB THE OTHER STUFF

//somehow join them together
func (m *postgresDBRepo) AddSongToPlaylist(song models.Song, playlist models.Playlist) (models.SongPlaylist, error) {
	query := "insert into songplaylist (playlist_id, song_id) values($1, $2, $3) returning *"
	var songplaylists models.SongPlaylist

	row := m.DB.QueryRow(query, playlist.Playlist_id, song.Song_id)

	err := row.Scan(query, &songplaylists.Playlist_id, &songplaylists.Song_id)
	if err != nil {
		log.Println(err)
	}

	return songplaylists, nil
}

//join song to album based on this
func (m *postgresDBRepo) AddSongToAlbum(res models.Song, album models.Album) (models.AlbumSong, error) {
	// query := "select song from song where title == $1"
	// add_query := "insert into song(album) values ($1)"

	var albumsong models.AlbumSong

	query := `insert into albumsong(album_id, song_id)
	values ($1, $2) returning *`

	row := m.DB.QueryRow(query, album.Album_id, res.Song_id)

	err := row.Scan(&albumsong.Name, &albumsong.Album_id, &albumsong.Song_id)//check for emptpy vals or errors in row
	if err != nil {
		log.Println(err)
	}

	return albumsong, nil
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

func(m *postgresDBRepo) UpdateUser(user models.Users) (models.Users, error){

	var users models.Users

	query := 
	`UPDATE Users
	SET (username, password, first_name, last_name, gender, admin) = ($1,$2,$3,$4,$5,$6)
    WHERE user_id = $7
			RETURNING *`

	row := m.DB.QueryRow(query, user.Username, user.Password, user.First_name, user.Last_name, user.Gender, user.Admin, user.User_id)

	err := row.Scan(&users.User_id, &users.Username, &user.Password, &user.First_name, &user.Last_name, &user.Gender, &user.Admin)

	if err != nil{
		log.Println(err)
	}

	return users, nil


}

func(m* postgresDBRepo) UpdateArtist(artist models.Artist) (models.Artist, error){

	var artists models.Artist

	query := 
	`UPDATE Artist
	SET (name, location) = ($1, $2)`

	row := m.DB.QueryRow(query, artist.Name, artist.Location)

	err := row.Scan(&artist.Name, & artist.Artist_id, &artist.Location, &artist.Join_date)

	if err != nil{
		log.Println(err)
	}

	return artists, nil
}

func(m* postgresDBRepo) UpdateSong(song models.Song) (models.Song, error){

	var songs models.Song

	query := `
	UPDATE SONG
	SET (title, duration, total_plays) = ($1, $2, 0)`

	row := m.DB.QueryRow(query, song.Title, song.Duration, song.Total_plays)

	err := row.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album, &song.Total_plays)

	if err != nil{
		log.Println(err)
	}

	return songs, nil

}


//delete from playlist/album - jonathan's job
//functions to add song to album/playlist and merge them together

