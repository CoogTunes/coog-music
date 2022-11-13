package dbrepo

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/CoogTunes/coog-music/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// USERS
func (m *postgresDBRepo) AddUser(res models.Users) (models.Users, error) {
	var user models.Users

	query := "insert into Users (username, password, first_name, last_name, admin_level) values ($1, $2, $3, $4, $5) RETURNING *"

	row := m.DB.QueryRow(query, res.Username, res.Password, res.First_name, res.Last_name, res.Admin_level)

	err := row.Scan(&user.User_id, &user.Username, &user.Password, &user.First_name, &user.Last_name, &user.Admin_level)
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

	err := rows.Scan(&user.User_id, &user.Username, &user.Password, &user.First_name, &user.Last_name, &user.Admin_level)
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

func (m *postgresDBRepo) GetArtistsAndSongs() ([]models.Artist, error) {
	var artists []models.Artist

	query := `SELECT a.name, a.artist_id, a.location, a.join_date, s.song_id, s.title, 
				s.release_date, s.duration, s.album_id, s.total_plays, s.album_id, al.name
				FROM Artist as a, Song as s, Album as al
				WHERE a.artist_id = s.artist_id AND al.album_id = s.album_id
				ORDER BY lower(a.name), a.artist_id, s.title`

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
		var song models.Song

		rows.Scan(&artist.Name, &artist.Artist_id, &artist.Location, &artist.Join_date, &song.Song_id, &song.Title,
			&song.Release_date, &song.Duration, &song.Album_id, &song.Total_plays, &song.Album_id, &song.Album)

		if err != nil {
			return nil, err
		}

		if len(artists) > 0 {
			if artist.Artist_id != artists[len(artists)-1].Artist_id {
				artists = append(artists, artist)
			}
		} else if len(artists) == 0 {
			artists = append(artists, artist)
		}

		song.Artist_id = artists[len(artists)-1].Artist_id
		song.Artist_name = artists[len(artists)-1].Name
		artists[len(artists)-1].Songs = append(artists[len(artists)-1].Songs, song)
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

func (m *postgresDBRepo) AddAlbum(album models.Album) (models.Album, error) {
	var addedAlbum models.Album

	query := "insert into Album (name, artist_id, date_added) values ($1, $2, to_date($3, 'YYYY-MM-DD')) RETURNING *"
	row := m.DB.QueryRow(query, album.Name, album.Artist_id, album.Date_added)

	err := row.Scan(&addedAlbum.Name, &addedAlbum.Artist_id, &addedAlbum.Album_id, &addedAlbum.Date_added)
	if err != nil {
		log.Println(err)
	}
	return addedAlbum, nil
}

func (m *postgresDBRepo) Authenticate(email string, password string) (models.Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var userInfo models.Users
	row := m.DB.QueryRowContext(ctx, "select * from users where username = $1", email)
	err := row.Scan(&userInfo.User_id, &userInfo.Username, &userInfo.Password, &userInfo.First_name, &userInfo.Last_name, &userInfo.Admin_level)
	if err != nil {
		return userInfo, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return userInfo, err
	} else if err != nil {
		return userInfo, err
	}

	return userInfo, nil
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

<<<<<<< Updated upstream
	rows, err := m.DB.Query(query)
	if err != nil {
=======
	query := "select song.*, playlist.name from song, songplaylist where song_id = songplaylist.song_id and songplaylist.playlist_id in (SELECT playlist_id from playlist where LOWER(name) LIKE LOWER(%$1%)"
	rows, err := m.DB.Query(query, playlist_name)
	if err != nil{
>>>>>>> Stashed changes
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

<<<<<<< Updated upstream
		rows.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album_id, &song.Total_plays)

		if err != nil {
			log.Println(err)
		}
=======
		songs = append(songs, song)
	}
	return songs, nil
}

func (m *postgresDBRepo) GetSongsFromArtist(artist_name string)([]models.Song, error){//get all artists
	var songs []models.Song

	query := "select Song.*, Artist.name as ArtistName from Song, Artist where song.artist_id in (SELECT artist_id from artist where LOWER(name) like LOWER('%$1%')) ORDER BY Song.title"
	rows, err := m.DB.Query(query, artist_name)
	if err != nil{
		return nil, err
	}

	defer func (rows *sql.Rows){
		err := rows.Close()
		if err != nil{
			log.Println(err)
		}
	}(rows)

	for rows.Next(){
		var song models.Song
		
		rows.Scan(&song.Song_id, &song.Title, &song.Album_id, &song.Artist_id, &song.SongPath, &song.CoverPath, &song.Uploaded_date, &song.Total_plays, &song.Total_likes, &song.Artist_name)

		if err != nil{
			return nil, err
		}

		songs = append(songs, song)
	}
	return songs, nil
}

func (m *postgresDBRepo) GetSongsFromAlbum(album_name string)([]models.Album, error){//return album, get artist name
	var songs []models.Album

	query := "select Song.*, Album.name as AlbumName from Song, Album where song.album_id in (SELECT album_id from album where LOWER(name) like LOWER('%$1%')) ORDER BY Song.title"
	rows, err := m.DB.Query(query, album_name)
	if err != nil{
		return nil, err
	}

	defer func (rows *sql.Rows){
		err := rows.Close()
		if err != nil{
			log.Println(err)
		}
	}(rows)

	for rows.Next(){
		var song models.Song
		rows.Scan(&song.Song_id, &song.Title, &song.Album_id, &song.Artist_id, &song.SongPath, &song.CoverPath, &song.Uploaded_date, &song.Total_plays, &song.Total_likes, &song.Album)

		if err != nil{
			return nil, err
		}

>>>>>>> Stashed changes
		songs = append(songs, song)
	}
	return songs, nil
}

func (m *postgresDBRepo) GetSong(songID string) (models.Song, error) {

	var song models.Song

	query := "select * from song where song_id = $1;"

	row := m.DB.QueryRow(query, songID)
	log.Println("row", row)
	log.Println(row.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album, &song.Total_plays))

	//maybe call update playcount
	return song, nil

}

func (m *postgresDBRepo) UpdateSongCount(songWithId models.Song) (models.Song, error) {
	var song models.Song

	query := "UPDATE Song SET total_plays = total_plays + 1 where song_id = $1 returning *"

	row := m.DB.QueryRow(query, songWithId.Song_id)

	row.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album, &song.Total_plays)

	return song, nil

}

//TODO: ADD LINKING TABLES AND USE THEM TO GRAB THE OTHER STUFF

// somehow join them together
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

// join song to album based on this
func (m *postgresDBRepo) AddSongToAlbum(res models.Song, album models.Album) (models.AlbumSong, error) {
	// query := "select song from song where title == $1"
	// add_query := "insert into song(album) values ($1)"

	var albumsong models.AlbumSong

	query := `insert into albumsong(album_id, song_id)
	values ($1, $2) returning *`

	row := m.DB.QueryRow(query, album.Album_id, res.Song_id)

	err := row.Scan(&albumsong.Name, &albumsong.Album_id, &albumsong.Song_id) //check for emptpy vals or errors in row
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

func (m *postgresDBRepo) UpdateUser(user models.Users) (models.Users, error) {

	var users models.Users
	log.Println("user", user)
	query :=
		`UPDATE Users
	SET (username, password, first_name, last_name, admin_level) = ($1,$2,$3,$4,$5)
    WHERE user_id = $6
			RETURNING *`

	row := m.DB.QueryRow(query, user.Username, user.Password, user.First_name, user.Last_name, user.Admin_level, user.User_id)

	err := row.Scan(&users.User_id, &users.Username, &users.Password, &users.First_name, &users.Last_name, &users.Admin_level)

	if err != nil {
		log.Println(err)
	}

	return users, nil

}

func (m *postgresDBRepo) UpdateArtist(artist models.Artist) (models.Artist, error) {

	var artistToReturn models.Artist

	query :=
		`UPDATE Artist
	SET (name, location) = ($1, $2)
	WHERE name = $3 RETURNING *`

	row := m.DB.QueryRow(query, artist.Name, artist.Location, artist.Name)

	err := row.Scan(&artistToReturn.Name, &artistToReturn.Artist_id, &artistToReturn.Location, &artistToReturn.Join_date)

	if err != nil {
		log.Println(err)
	}

	return artistToReturn, nil
}

func (m *postgresDBRepo) UpdateSong(song models.Song) (models.Song, error) {

	var songToReturn models.Song

	query := `
	UPDATE SONG
	SET (title, duration) = ($1, $2) 
	WHERE song_id = $3 RETURNING *`

	row := m.DB.QueryRow(query, song.Title, song.Duration, song.Song_id)

	err := row.Scan(&songToReturn.Song_id, &songToReturn.Title, &songToReturn.Artist_id, &songToReturn.Release_date, &songToReturn.Duration, &songToReturn.Album_id, &songToReturn.Total_plays)

	if err != nil {
		log.Println(err)
	}

	return songToReturn, nil

}

func (m *postgresDBRepo) Follow(artistId int, userId int) (models.Followers, error) {

	var followers models.Followers

	query := `INSERT INTO FOLLOWERS (artist_id, user_id) values($1, $2) RETURNING *`

	row := m.DB.QueryRow(query, artistId, userId)

	err := row.Scan(&followers.Artist_id, &followers.User_id)

	if err != nil {
		log.Println(err)
	}

	return followers, nil
}

func (m *postgresDBRepo) RemoveUser(user_id int) error {

	query := `DELETE FROM USERS WHERE user_id = $1`

	_, err := m.DB.Exec(query, user_id)

	if err != nil {
		log.Println(err)
	}

	return nil
}

func (m *postgresDBRepo) RemoveSong(song_id int) error {

	query := `DELETE FROM SONG WHERE song_id = $1`

	_, err := m.DB.Exec(query, song_id)

	if err != nil {
		log.Println(err)
	}

	return nil
}

func (m *postgresDBRepo) RemovePlaylist(playlist_id int) error {

	query := `DELETE FROM PLAYLIST WHERE playlist_id = $1`

	_, err := m.DB.Exec(query, playlist_id)

	if err != nil {
		log.Println(err)
	}

	return nil
}

func (m *postgresDBRepo) RemoveAlbum(album_id int) error {

	query := `DELETE FROM ALBUM WHERE album_id = $1`

	_, err := m.DB.Exec(query, album_id)

	if err != nil {
		log.Println(err)
	}

	return nil
	}

func (m *postgresDBRepo) RemoveArtist(artist_id int) error {

	query := `DELETE FROM ARTIST WHERE artist_id = $1`

	_, err := m.DB.Exec(query, artist_id)

	if err != nil {
		log.Println(err)
	}

	return nil
}

func (m *postgresDBRepo) RemoveSongFromAlbum(song_id int, album_id int) error {

	query := "DELETE FROM ALBUMSONG WHERE song_id = $1 AND album_id = $2"

	_, err := m.DB.Exec(query, song_id, album_id)

	if err != nil {
		log.Println(err)
	}

	return nil

}

func (m *postgresDBRepo) RemoveSongFromPlaylist(song_id int, playlist_id int) error {

	query := "DELETE FROM SONGPLAYLIST WHERE song_id = $1 AND playlist_id = $2"

	_, err := m.DB.Exec(query, song_id, playlist_id)

	if err != nil {
		log.Println(err)
	}

	return nil
}

func(m *postgresDBRepo) GetNumberOfUsers()(models.Users, error){

	var user models.Users

	query := "SELECT MAX(user_id) FROM USERS WHERE admin_level = 1 RETURNING user_id"
	row := m.DB.QueryRow(query)

	err := row.Scan(&user.User_id)
	if err != nil{
		log.Println(err)
	}

	return user, nil
}

func(m *postgresDBRepo) GetNumberOfArtists()(models.Users, error){

	var user models.Users

	query := "SELECT MAX(user_id) FROM USERS WHERE admin_level = 2 RETURNING user_id"
	row := m.DB.QueryRow(query)

	err := row.Scan(&user.User_id)
	if err != nil{
		log.Println(err)
	}

	return user, nil

}

func(m *postgresDBRepo) GetNumberOfSongs()(models.Song, error){
	var song models.Song

	query := "SELECT MAX(song_id) FROM SONG"
	row := m.DB.QueryRow(query)

	err := row.Scan(&song.Song_id)
	if err != nil{
		log.Println(err)
	}
	return song, nil
}

func(m *postgresDBRepo) GetNumberOfPlaylists()(models.Playlist, error){
	var playlist models.Playlist

	query := "SELECT MAX(playlist_id) FROM PLAYLIST"

	row := m.DB.QueryRow(query)

	err := row.Scan(&playlist.Playlist_id)
	if err != nil{
	log.Println(err)
	}
	return playlist, nil
}

<<<<<<< Updated upstream
func (m *postgresDBRepo) GetSongsFromPlaylist (playlist_id int)([]models.Song, error){
	var songs []models.Song

	query := "select song from songplaylist, song where songplaylist.playlist_id = $1"
	rows, err := m.DB.Query(query, playlist_id)
	if err != nil{
		return nil, err
	}

	defer func (rows *sql.Rows){
		err := rows.Close()
		if err != nil{
			log.Println(err)
		}
	}(rows)

	for rows.Next(){
		var song models.Song
		
		rows.Scan(&song.Song_id, &song.Artist_id)

		if err != nil{
			return nil, err
		}

		songs = append(songs, song)
	}
	return songs, nil
}
=======
func(m *postgresDBRepo) AddLikeToSong(song_id int)error{
	query := "UPDATE Song SET total_likes = total_likes + 1 where song_id = $1"
>>>>>>> Stashed changes

	_, err := m.DB.Exec(query, song_id)

	if err != nil{
		log.Println(err)
	}
	return nil
}
//delete from playlist/album
//delete artist, user, song
//functions to add song to album/playlist and merge them together
