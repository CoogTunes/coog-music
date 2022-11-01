package models

import (
	"time"
)

type Users struct {
	User_id    string
	Username   string
	Password   string
	First_name string
	Last_name  string
	Gender     string
	Admin      bool
}

type Song struct {
	Song_id      int
	Title        string
	Artist_id    string
	Release_date string
	Duration     string
	Artist_name  string
	Album        string
	Total_plays  int
}

type Artist struct {
	Name      string
	Artist_id int
	Location  string
	Join_date string
	Songs     []int
	Admin     bool
}

type Playlist struct {
	Song_id            int
	Name            string
	Playlist_length int
	Playlist_id     int
	User_id         int
}

type Album struct {
	Name         string
	Album_id     int
	Artist_id    int
	Date_added   time.Time
	Song_id      int
}

type Songplay struct {
	Songplay_id int
	Session_id  int
	Level       string
	Song_id     int
	Artist_id   int
}
