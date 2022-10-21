package models

type Users struct {
	Username   string
	Password   string
	First_name string
	Last_name  string
	Gender     string
}

type Song struct{
	Title string
	Artist_name string
	Album string
}

type Artist struct{
	Name string
	Publisher string
}
