package models

import (
	"time"
)

type Users struct {
	Username   string
	Password   string
	First_name string
	Last_name  string
	Gender     string
}

type Song struct {
	Song_id     int
	Title       string
	Artist_name string
	Album       string
}

type Artist struct {
	Name      string
	Artist_id int64
	Location  string
	Join_date string
	Songs     []int
	Admin     bool
	Publisher string
}

type Playlist struct {
	Songs           int
	Playlist_length int
	Playlist_id     int
}

type Album struct {
	Name         string
	Album_id     int
	Date_added   time.Time
	Publisher_id int
}

type Songplay struct {
	Songplay_id int
	Session_id  int
	Level       string
	Song_id     int
	Artist_id   int
}

type SongPlaylist struct {
	Song_id     int
	Playlist_id int
}

type AlbumSong struct {
	Album_id int
	Song_id  int
	Name     string
}
