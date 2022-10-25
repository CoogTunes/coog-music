package models

// type Artist struct {
// 	name      string
// 	artist_id int64
// 	location  string
// 	join_date string
// 	songs     []int
// 	admin     bool
// }

type Artist struct {
	// 	Description string `json:"description"`
	Name      string `json:"name"`
	Artist_id int64  `json:"artist_id"`
	Location  string `json:"location"`
	Join_date string `json:"join_date"`
	Songs     []int  `json:"songs"`
	Admin     bool   `json:"admin"`
}
