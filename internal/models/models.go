package models

type Users struct {
	User_id     int
	Username    string
	Password    string
	First_name  string
	Last_name   string
	Admin_level int
	JoinedDate  string
}

type Song struct {
	Song_id       int
	Title         string
	Artist_id     int
	SongPath      string
	CoverPath     string
	Uploaded_date string
	Album         string
	Album_id      int
	Total_plays   int
	Total_likes   int
	Artist_name   string
}

type Artist struct {
	Name      string
	Artist_id int
	Location  string
	Join_date string
	Songs     []Song
}

type Playlist struct {
	User_id     int
	Playlist_id int
	Name        string
}

type Album struct {
	Name       string
	Album_id   int
	Artist_id  int
	Song_names []string
	Date_added string
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

type Followers struct {
	User_id   int
	Artist_id int
}

type LikesReport struct {
	Likes         int
	Dislikes      int
	Song_id       int
	Song_title    string
	Artist_name   string
	Album_name    string
	Uploaded_date string
}
