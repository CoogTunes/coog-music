package models

type Artist struct {
	Name      string
	Artist_id int64
	Location  string
	Join_date string
	Songs     []int
	Admin     bool
}
