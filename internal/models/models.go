package models

type Users struct {
	User_id           int
	Username          string
	Password          string
	First_name        string
	Last_name         string
	Concatenated_name string
	Admin_level       int
	JoinedDate        string
}

type Song struct {
	Likes         int
	Dislikes      int
	Song_id       int
	Title         string
	Artist_id     int
	SongPath      string
	CoverPath     string
	Uploaded_date string
	Album         string
	Album_id      int
	Total_plays   int
	Artist_name   string
	Duration      string
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

// type LikesReport struct {      // Just using song struct instead
// 	Likes         int
// 	Dislikes      int
// 	Song_id       int
// 	Song_title    string
// 	Artist_name   string
// 	Album_name    string
// 	Uploaded_date string
// 	Song_path     string
// 	Cover_path    string
// 	Artist_id     int
// 	Album_id      int
// }

type UserReport struct {
	User_id           int
	Username          string
	First_name        string
	Last_name         string
	Admin_level       int
	JoinedDate        string
	Playlist_count    int
	Liked_songs_count int
	Common_artist     string
}

type ArtistReport struct {
	Name            string
	Artist_id       int
	Join_date       string
	Num_songs       int
	Num_Albums      int
	Total_Plays     int
	Avg_Plays       int
	Most_liked_song string
}

type DisplaySongInfo struct {
	SongPath     string
	CoverPath    string
	Title        string
	SongID       int
	Album        string
	Artist       string
	Likes        int
	Dislikes     int
	UploadedDate string
	ArtistID     int
	AlbumID      int
}

type Messages struct {
	User_id      int
	Message      string
	Created_date string
	IsRead       bool
	Message_id   int
}
