package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
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

	err := row.Scan(&user.User_id, &user.Username, &user.Password, &user.First_name, &user.Last_name, &user.Admin_level, &user.JoinedDate)
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
func (m *postgresDBRepo) AddArtist(res models.Artist) error {

	query := "insert into Artist (name, artist_id, location, join_date) values ($1, $2, $3, to_date($4, 'YYYY-MM-DD'))"
	_, err := m.DB.Exec(query, res.Name, res.Artist_id, res.Location, time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *postgresDBRepo) AddPlaylist(res models.Playlist) (models.Playlist, error) {
	var playlist models.Playlist
	query := "insert into Playlist (user_id, name) values($1, $2) RETURNING *"
	row := m.DB.QueryRow(query, res.User_id, res.Name)
	err := row.Scan(&playlist.User_id, &playlist.Name, &playlist.Playlist_id)
	if err != nil {
		log.Println(err)
		return playlist, err
	}

	return playlist, nil

}

func (m *postgresDBRepo) AddPlaylistSong(songID int, playlistID int) error {
	query := "insert into SongPlaylist (song_id, playlist_id) values ($1, $2)"
	_, err := m.DB.Exec(query, songID, playlistID)
	if err != nil {
		log.Println("Cannot add song playlist")
		log.Println(err)
		return err
	}
	return nil
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

//func (m *postgresDBRepo) GetArtistsAndSongs() ([]models.Artist, error) {
//	var artists []models.Artist
//
//	query := `SELECT a.name, a.artist_id, a.location, a.join_date, s.song_id, s.title,
//				s.release_date, s.duration, s.album_id, s.total_plays, s.album_id, al.name
//				FROM Artist as a, Song as s, Album as al
//				WHERE a.artist_id = s.artist_id AND al.album_id = s.album_id
//				ORDER BY lower(a.name), a.artist_id, s.title`
//
//	rows, err := m.DB.Query(query)
//	if err != nil {
//		return nil, err
//	}
//
//	defer func(rows *sql.Rows) {
//		err := rows.Close()
//		if err != nil {
//			log.Println(err)
//		}
//	}(rows)
//
//	for rows.Next() {
//		var artist models.Artist
//		var song models.Song
//
//		rows.Scan(&artist.Name, &artist.Artist_id, &artist.Location, &artist.Join_date, &song.Song_id, &song.Title,
//			&song.Release_date, &song.Duration, &song.Album_id, &song.Total_plays, &song.Album_id, &song.Album)
//
//		if err != nil {
//			return nil, err
//		}
//
//		if len(artists) > 0 {
//			if artist.Artist_id != artists[len(artists)-1].Artist_id {
//				artists = append(artists, artist)
//			}
//		} else if len(artists) == 0 {
//			artists = append(artists, artist)
//		}
//
//		song.Artist_id = artists[len(artists)-1].Artist_id
//		song.Artist_name = artists[len(artists)-1].Name
//		artists[len(artists)-1].Songs = append(artists[len(artists)-1].Songs, song)
//	}
//	return artists, nil
//}

//func (m *postgresDBRepo) GetArtistName(artist_id int) (str string, err error) {
//	var song models.Song
//	query := `SELECT name FROM Artist as A, Song as S WHERE A.artist_id = S.artist_id AND A.artist_id = $1`
//	row := m.DB.QueryRow(query, artist_id)
//	err2 := row.Scan(&song.Artist_name)
//	if err2 != nil {
//		log.Println(err)
//	}
//
//	return song.Artist_name, err
//}

func (m *postgresDBRepo) AddAlbum(album models.Album) (models.Album, error) {
	var addedAlbum models.Album

	query := "insert into Album (name, artist_id, date_added) values ($1, $2, to_date($3, 'YYYY-MM-DD')) RETURNING *"
	row := m.DB.QueryRow(query, album.Name, album.Artist_id, time.Now())

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
	err := row.Scan(&userInfo.User_id, &userInfo.Username, &userInfo.Password, &userInfo.First_name, &userInfo.Last_name, &userInfo.Admin_level, &userInfo.JoinedDate)
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

// SONG
func (m *postgresDBRepo) AddSong(res models.Song) error {
	query := "insert into song (title, artist_id, song_path, cover_path, uploaded_date, duration) values ($1, $2, $3, $4, to_date($5, 'YYYY-MM-DD'), $6)"

	_, err := m.DB.Exec(query, res.Title, res.Artist_id, res.SongPath, res.CoverPath, res.Uploaded_date, res.Duration)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// SONGS
func (m *postgresDBRepo) AddSongForAlbum(res models.Song) error {
	query := "insert into song (title, artist_id, album_id, song_path, cover_path, uploaded_date,duration) values ($1, $2, $3, $4, $5, to_date($6, 'YYYY-MM-DD'), $7)"

	_, err := m.DB.Exec(query, res.Title, res.Artist_id, res.Album_id, res.SongPath, res.CoverPath, res.Uploaded_date, res.Duration)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetSongsForLikePage(userId int) ([]models.Song, error) {
	var songs []models.Song
	query := `select * from likes_view 
					where likes_view.song_id in (select likes.song_id from likes where user_id = $1 and likes.islike = true) 
					order by artist_name, album_name, title;`

	rows, err := m.DB.Query(query, userId)
	if err != nil {
		log.Println(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var currentRow models.Song

		err := rows.Scan(&currentRow.Likes, &currentRow.Dislikes, &currentRow.Song_id, &currentRow.Title, &currentRow.Album_id, &currentRow.Artist_id,
			&currentRow.SongPath, &currentRow.CoverPath, &currentRow.Uploaded_date, &currentRow.Total_plays, &currentRow.Duration, &currentRow.Artist_name, &currentRow.Album)

		if err != nil {
			return nil, err
		}
		songs = append(songs, currentRow)
	}
	return songs, err
}

func (m *postgresDBRepo) GetSongsFromPlaylist(userId int, playlist_id int) ([]models.Song, error) {
	var songsInfo []models.Song

	query := `select coalesce(likes.isLike, false), likes_view.* from likes_view
			left join likes on likes.song_id = likes_view.song_id and likes.user_id = $1
			where likes_view.song_id in (select songplaylist.song_id from songplaylist where songplaylist.playlist_id = $2)`
	rows, err := m.DB.Query(query, userId, playlist_id)
	if err != nil {
		log.Println("Cannot execute query")
		return songsInfo, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var currentRow models.Song
		rows.Scan(&currentRow.IsLike, &currentRow.Likes, &currentRow.Dislikes, &currentRow.Song_id, &currentRow.Title, &currentRow.Album_id, &currentRow.Artist_id,
			&currentRow.SongPath, &currentRow.CoverPath, &currentRow.Uploaded_date, &currentRow.Total_plays, &currentRow.Duration, &currentRow.Artist_name, &currentRow.Album)

		if err != nil {
			log.Println("Cannot scan row")
			return songsInfo, err
		}
		songsInfo = append(songsInfo, currentRow)
	}
	return songsInfo, nil
}

func (m *postgresDBRepo) GetSongsFromArtist(artist_name string) (map[string][]models.Song, error) {
	// var songsInfo []models.Song
	songsInfo := make(map[string][]models.Song)
	fmt.Println(artist_name, "artist name")
	query := "select * from likes_view where LOWER(likes_view.artist_name) LIKE LOWER('" + artist_name + "%')"
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
		var currentRow models.Song
		rows.Scan(&currentRow.Likes, &currentRow.Dislikes, &currentRow.Song_id, &currentRow.Title, &currentRow.Album_id, &currentRow.Artist_id,
			&currentRow.SongPath, &currentRow.CoverPath, &currentRow.Uploaded_date, &currentRow.Total_plays, &currentRow.Duration, &currentRow.Artist_name, &currentRow.Album)

		if err != nil {
			log.Println("Cannot scan row")
			return songsInfo, err
		}

		songsInfo[currentRow.Album] = append(songsInfo[currentRow.Album], currentRow)
	}
	return songsInfo, nil
}

func (m *postgresDBRepo) GetSongsFromArtistByID(userId int, artist_name string, artistID int) (map[string][]models.Song, error) {
	songsInfo := make(map[string][]models.Song)
	// query := "select * from likes_view where LOWER(likes_view.artist_name) LIKE LOWER('" + artist_name + "%') and likes_view.artist_id = $1"
	query1 := "select coalesce(likes.islike, false), likes_view.* from likes_view "
	query2 := "left join likes on likes.song_id = likes_view.song_id and likes.user_id = $1 "
	query3 := "where LOWER(likes_view.artist_name) LIKE LOWER('" + artist_name + "%') and likes_view.artist_id = $2 order by likes_view.title"
	query := query1 + query2 + query3

	rows, err := m.DB.Query(query, userId, artistID)
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
		var currentRow models.Song
		rows.Scan(&currentRow.IsLike, &currentRow.Likes, &currentRow.Dislikes, &currentRow.Song_id, &currentRow.Title, &currentRow.Album_id, &currentRow.Artist_id,
			&currentRow.SongPath, &currentRow.CoverPath, &currentRow.Uploaded_date, &currentRow.Total_plays, &currentRow.Duration, &currentRow.Artist_name, &currentRow.Album)

		if err != nil {
			log.Println("Cannot scan row")
			return songsInfo, err
		}
		songsInfo[currentRow.Album] = append(songsInfo[currentRow.Album], currentRow)
	}
	return songsInfo, nil
}

func (m *postgresDBRepo) GetSongsFromAlbum(album_name string) ([]models.Song, error) {
	var songsInfo []models.Song

	query := "select * from likes_view where LOWER(likes_view.album_name) LIKE LOWER ('" + album_name + "%')"
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
		var currentRow models.Song
		rows.Scan(&currentRow.Likes, &currentRow.Dislikes, &currentRow.Song_id, &currentRow.Title, &currentRow.Album_id, &currentRow.Artist_id,
			&currentRow.SongPath, &currentRow.CoverPath, &currentRow.Uploaded_date, &currentRow.Total_plays, &currentRow.Duration, &currentRow.Artist_name, &currentRow.Album)

		if err != nil {
			log.Println("Cannot scan row")
			return songsInfo, err
		}
		songsInfo = append(songsInfo, currentRow)
	}
	return songsInfo, nil
}

func (m *postgresDBRepo) GetSongsFromAlbumByID(album_name string, albumID int) ([]models.Song, error) {
	var songsInfo []models.Song

	query := "select * from likes_view where LOWER(likes_view.album_name) LIKE LOWER ('" + album_name + "%') and likes_view.album_id = $1"
	rows, err := m.DB.Query(query, albumID)
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
		var currentRow models.Song
		rows.Scan(&currentRow.Likes, &currentRow.Dislikes, &currentRow.Song_id, &currentRow.Title, &currentRow.Album_id, &currentRow.Artist_id,
			&currentRow.SongPath, &currentRow.CoverPath, &currentRow.Uploaded_date, &currentRow.Total_plays, &currentRow.Duration, &currentRow.Artist_name, &currentRow.Album)

		if err != nil {
			log.Println("Cannot scan row")
			return songsInfo, err
		}
		songsInfo = append(songsInfo, currentRow)
	}
	return songsInfo, nil
}

func (m *postgresDBRepo) GetSongsByName(song_name string) ([]models.Song, error) {
	var songs []models.Song

	// query := `select al.name, al.album_id, s.artist_id, s.title, s.song_id, s.cover_path, s.song_path, s.uploaded_date, ar.name
	// from album as al, artist as ar, song as s
	// where s.title = $1 and s.artist_id = al.artist_id and s.artist_id = ar.artist_id`

	// query := "select al.name, al.album_id, s.artist_id, s.title, s.song_id, s.cover_path, s.song_path, s.uploaded_date, ar.name from album as al, artist as ar, song as s where LOWER(s.title) LIKE LOWER ('" + song_name + "%') and s.artist_id = al.artist_id and s.artist_id = ar.artist_id"

	query := "select * from likes_view where  LOWER(likes_view.title) LIKE LOWER ('" + song_name + "%')  FETCH FIRST 7 ROWS ONLY"
	rows, err := m.DB.Query(query)
	if err != nil {
		log.Println("Cannot get any rows")
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var currentRow models.Song

		rows.Scan(&currentRow.Likes, &currentRow.Dislikes, &currentRow.Song_id, &currentRow.Title, &currentRow.Album_id, &currentRow.Artist_id,
			&currentRow.SongPath, &currentRow.CoverPath, &currentRow.Uploaded_date, &currentRow.Total_plays, &currentRow.Duration, &currentRow.Artist_name, &currentRow.Album)

		if err != nil {
			log.Println("Cannot scan any rows")
			return nil, err
		}

		songs = append(songs, currentRow)
	}

	return songs, nil
}

func (m *postgresDBRepo) GetTopSongs() ([]models.Song, error) {
	var songs []models.Song
	query := "select * from likes_view order by total_plays desc, title Fetch first 14 rows only"
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
		var currentRow models.Song
		rows.Scan(&currentRow.Likes, &currentRow.Dislikes, &currentRow.Song_id, &currentRow.Title, &currentRow.Album_id, &currentRow.Artist_id,
			&currentRow.SongPath, &currentRow.CoverPath, &currentRow.Uploaded_date, &currentRow.Total_plays, &currentRow.Duration, &currentRow.Artist_name, &currentRow.Album)
		log.Println(currentRow.SongPath)
		if err != nil {
			return nil, err
		}

		songs = append(songs, currentRow)
	}
	return songs, nil

}

//func (m *postgresDBRepo) GetSongs() ([]models.Song, error) {
//	var songs []models.Song
//	// probably need to add a where statement and get rid of *
//	query := "SELECT * FROM song"
//
//	rows, err := m.DB.Query(query)
//	if err != nil {
//		return nil, err
//	}
//
//	defer func(rows *sql.Rows) {
//		err := rows.Close()
//		if err != nil {
//			log.Println(err)
//		}
//	}(rows)
//
//	for rows.Next() {
//		var song models.Song
//
//		rows.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album_id, &song.Total_plays)
//
//		if err != nil {
//			log.Println(err)
//		}
//		songs = append(songs, song)
//	}
//	return songs, nil
//}

//func (m *postgresDBRepo) GetSong(songID string) (models.Song, error) {
//
//	var song models.Song
//
//	query := "select * from song where song_id = $1"
//
//	row := m.DB.QueryRow(query, songID)
//	log.Println("row", row)
//	log.Println(row.Scan(&song.Song_id, &song.Title, &song.Artist_id, &song.Release_date, &song.Duration, &song.Album, &song.Total_plays))
//
//	//maybe call update playcount
//	return song, nil
//
//}

func (m *postgresDBRepo) UpdateSongCount(id int) (models.Song, error) {
	var song models.Song

	query := "UPDATE Song SET total_plays = total_plays + 1 where song_id = $1 returning *"

	row := m.DB.QueryRow(query, id)

	row.Scan(&song.Song_id, &song.Title, &song.Album_id, &song.Artist_id, &song.SongPath, &song.CoverPath, &song.Uploaded_date, &song.Total_plays, &song.Duration)

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

	query := `insert into albumsong(album_id, song_id) values ($1, $2) returning *`

	row := m.DB.QueryRow(query, album.Album_id, res.Song_id)

	err := row.Scan(&albumsong.Name, &albumsong.Album_id, &albumsong.Song_id) //check for emptpy vals or errors in row
	if err != nil {
		log.Println(err)
	}

	return albumsong, nil
}

func (m *postgresDBRepo) GetPlaylists(userId int) ([]models.Playlist, error) {
	var playlists []models.Playlist

	query := "SELECT name, playlist_id FROM PLAYLIST where user_id = $1"

	rows, err := m.DB.Query(query, userId)
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

		rows.Scan(&playlist.Name, &playlist.Playlist_id)
		playlist.User_id = userId
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

// func (m *postgresDBRepo) UpdateSong(song models.Song) (models.Song, error) {

// 	var songToReturn models.Song

// 	query := `
// 	UPDATE SONG
// 	SET (title, duration) = ($1, $2)
// 	WHERE song_id = $3 RETURNING *`

// 	row := m.DB.QueryRow(query, song.Title, song.Song_id)

// 	err := row.Scan(&songToReturn.Song_id, &songToReturn.Title, &songToReturn.Artist_id, &songToReturn.Album_id, &songToReturn.Total_plays)

// 	if err != nil {
// 		log.Println(err)
// 	}

// 	return songToReturn, nil

// }

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

func (m *postgresDBRepo) GetNumberOfUsers() (models.Users, error) {

	var user models.Users

	query := "SELECT MAX(user_id) FROM USERS WHERE admin_level = 1 RETURNING user_id"
	row := m.DB.QueryRow(query)

	err := row.Scan(&user.User_id)
	if err != nil {
		log.Println(err)
	}

	return user, nil
}

func (m *postgresDBRepo) GetNumberOfArtists() (models.Users, error) {

	var user models.Users

	query := "SELECT MAX(user_id) FROM USERS WHERE admin_level = 2 RETURNING user_id"
	row := m.DB.QueryRow(query)

	err := row.Scan(&user.User_id)
	if err != nil {
		log.Println(err)
	}

	return user, nil

}

func (m *postgresDBRepo) GetNumberOfSongs() (models.Song, error) {
	var song models.Song

	query := "SELECT MAX(song_id) FROM SONG"
	row := m.DB.QueryRow(query)

	err := row.Scan(&song.Song_id)
	if err != nil {
		log.Println(err)
	}
	return song, nil
}

func (m *postgresDBRepo) GetNumberOfPlaylists() (models.Playlist, error) {
	var playlist models.Playlist

	query := "SELECT MAX(playlist_id) FROM PLAYLIST"

	row := m.DB.QueryRow(query)

	err := row.Scan(&playlist.Playlist_id)
	if err != nil {
		log.Println(err)
	}
	return playlist, nil
}

func (m *postgresDBRepo) AddOrUpdateLikeValue(islike bool, songId int, userId int) error {
	query := "insert into likes (islike,song_id,user_id) values ($1,$2,$3) RETURNING *"
	_, err := m.DB.Exec(query, islike, songId, userId)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (m *postgresDBRepo) SendUpdatedLikeValue(songID int) (models.Song, error) {
	var currentRow models.Song
	query := "select * from likes_view where likes_view.song_id = $1"
	row := m.DB.QueryRow(query, songID)
	err := row.Scan(&currentRow.Likes, &currentRow.Dislikes, &currentRow.Song_id, &currentRow.Title, &currentRow.Album_id, &currentRow.Artist_id,
		&currentRow.SongPath, &currentRow.CoverPath, &currentRow.Uploaded_date, &currentRow.Total_plays, &currentRow.Duration, &currentRow.Artist_name, &currentRow.Album)
	if err != nil {
		log.Println(err)
		return currentRow, err
	}
	return currentRow, nil
}

// REPORTS
func (m *postgresDBRepo) GetLikesReport(minLikes int, maxLikes int, minDate string, maxDate string) ([]models.Song, error) {
	var likesReport []models.Song
	query := `select * from likes_view 
				where likes_view.likes >= $1
				AND likes_view.likes <= $2
				AND likes_view.uploaded_date >= $3
				AND likes_view.uploaded_date <= $4
				ORDER BY likes_view.likes DESC, likes_view.title`

	rows, err := m.DB.Query(query, minLikes, maxLikes, minDate, maxDate)
	if err != nil {
		log.Println(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var currentRow models.Song

		err := rows.Scan(&currentRow.Likes, &currentRow.Dislikes, &currentRow.Song_id, &currentRow.Title, &currentRow.Album_id, &currentRow.Artist_id,
			&currentRow.SongPath, &currentRow.CoverPath, &currentRow.Uploaded_date, &currentRow.Total_plays, &currentRow.Duration, &currentRow.Artist_name, &currentRow.Album)

		if err != nil {
			return nil, err
		}
		likesReport = append(likesReport, currentRow)
	}
	return likesReport, nil
}

func (m *postgresDBRepo) GetInitialUsersReport() ([]models.UserReport, error) {
	var usersReport []models.UserReport
	query := `select * from usersreport where admin_level = 1 or admin_level = 2 ORDER BY join_date, last_name FETCH FIRST 14 ROWS ONLY`
	rows, err := m.DB.Query(query)
	if err != nil {
		log.Println(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var currentRow models.UserReport

		err := rows.Scan(&currentRow.User_id, &currentRow.Username, &currentRow.First_name, &currentRow.Last_name, &currentRow.Admin_level, &currentRow.JoinedDate, &currentRow.Playlist_count,
			&currentRow.Liked_songs_count, &currentRow.Common_artist)

		if err != nil {
			return nil, err
		}
		usersReport = append(usersReport, currentRow)
	}
	return usersReport, nil
}

func (m *postgresDBRepo) GetUsersReport(minDate string, maxDate string) ([]models.UserReport, error) {
	var usersReport []models.UserReport
	query := `select * from usersreport where (join_date >= to_date($1, 'YYYY-MM-DD') AND join_date <= to_date($2, 'YYYY-MM-DD')) and (admin_level = 1 or admin_level = 2) ORDER BY join_date, last_name`

	rows, err := m.DB.Query(query, minDate, maxDate)
	if err != nil {
		log.Println(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var currentRow models.UserReport

		err := rows.Scan(&currentRow.User_id, &currentRow.Username, &currentRow.First_name, &currentRow.Last_name, &currentRow.Admin_level, &currentRow.JoinedDate, &currentRow.Playlist_count,
			&currentRow.Liked_songs_count, &currentRow.Common_artist)

		if err != nil {
			return nil, err
		}
		usersReport = append(usersReport, currentRow)
	}
	return usersReport, nil
}

func (m *postgresDBRepo) GetArtistReport(minDate string, maxDate string) ([]models.ArtistReport, error) {
	var artistReport []models.ArtistReport
	query := `select * from artistsReport where join_date >= to_date($1, 'YYYY-MM-DD') AND join_date <= to_date($2, 'YYYY-MM-DD') ORDER BY join_date desc, name`

	rows, err := m.DB.Query(query, minDate, maxDate)
	if err != nil {
		log.Println(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var currentRow models.ArtistReport

		err := rows.Scan(&currentRow.Name, &currentRow.Artist_id, &currentRow.Join_date, &currentRow.Num_songs, &currentRow.Num_Albums, &currentRow.Total_Plays, &currentRow.Avg_Plays, &currentRow.Most_liked_song)

		if err != nil {
			return nil, err
		}
		artistReport = append(artistReport, currentRow)
	}
	return artistReport, nil
}

func (m *postgresDBRepo) GetSongReport(minDate string, maxDate string, min_plays int, max_plays int) ([]models.Song, error) {
	var songReport []models.Song

	query := `select * from likes_view
				where likes_view.uploaded_date >= $1
				AND likes_view.uploaded_date <= $2
				AND likes_view.total_plays >= $3
				AND likes_view.total_plays <= $4
				ORDER BY likes_view.total_plays DESC, likes_view.title`

	rows, err := m.DB.Query(query, minDate, maxDate, min_plays, max_plays)
	if err != nil {
		log.Println(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var currentRow models.Song
		err := rows.Scan(&currentRow.Likes, &currentRow.Dislikes, &currentRow.Song_id, &currentRow.Title, &currentRow.Album_id, &currentRow.Artist_id,
			&currentRow.SongPath, &currentRow.CoverPath, &currentRow.Uploaded_date, &currentRow.Total_plays, &currentRow.Duration, &currentRow.Artist_name, &currentRow.Album)

		if err != nil {
			return nil, err
		}
		songReport = append(songReport, currentRow)
	}
	return songReport, nil
}

func (m *postgresDBRepo) UpdateMessages(user_id int) ([]models.Messages, error) {
	var messages []models.Messages

	checkQuery := `UPDATE MESSAGES SET isRead = TRUE WHERE user_id = $1 and isRead = FALSE RETURNING *`
	rows, err := m.DB.Query(checkQuery, user_id)
	if err != nil {
		log.Println(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var currentMessage models.Messages
		err := rows.Scan(&currentMessage.User_id, &currentMessage.Message, &currentMessage.Created_date, &currentMessage.IsRead, &currentMessage.Message_id)

		if err != nil {
			return nil, err
		}
		messages = append(messages, currentMessage)
	}
	return messages, nil

}

// func (m *postgresDBRepo) UpdateMessage(user_id int) error {
// 	updateQuery := "UPDATE MESSAGES SET isRead = TRUE WHERE user_id = $1 and isRead = FALSE"

// 	_, err := m.DB.Exec(updateQuery, user_id)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return nil
// }

//add like function to increment in song
//search function for album, artist, playlist, get every song under these params
