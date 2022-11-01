package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/DeLuci/coog-music/internal/config"
	"github.com/DeLuci/coog-music/internal/driver"
	"github.com/DeLuci/coog-music/internal/models"
	"github.com/go-chi/chi/v5"

	"github.com/DeLuci/coog-music/internal/repository"
	"github.com/DeLuci/coog-music/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// USERS
func (m *Repository) AddUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// get fields
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	first_name := r.Form.Get("first_name")
	last_name := r.Form.Get("last_name")
	gender := r.Form.Get("gender")
	stringAdmin := r.Form.Get("admin")

	// make sure admin is a boolean
	boolAdmin, err := strconv.ParseBool(stringAdmin)
	if err != nil {
		log.Println(err)
	}

	// -1 is just to pass the user object to postgres.go, but it will not be used.
	newUser := models.Users{User_id: -1, Username: username, Password: password, First_name: first_name, Last_name: last_name, Gender: gender, Admin: boolAdmin}

	addedUser, err := m.DB.AddUser(newUser)

	returnAsJSON(addedUser, w, err)
}

func (m *Repository) GetUser(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	user, err := m.DB.GetUser(id)

	returnAsJSON(user, w, err)
}

// ARTISTS
func (m *Repository) AddArtist(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// get fields
	name := r.Form.Get("name")
	artist_id := r.Form.Get("artist_id")
	location := r.Form.Get("location")
	join_date := r.Form.Get("join_date")
	var songs []int

	int_artist_id, err := strconv.Atoi(artist_id)
	if err != nil {
		log.Println(err)
	}
	// joindate and songs[] should be empty to start.
	artistToAdd := models.Artist{Name: name, Artist_id: int_artist_id, Location: location, Join_date: join_date, Songs: songs}

	addedArtist, err := m.DB.AddArtist(artistToAdd)

	returnAsJSON(addedArtist, w, err)

}
func (m *Repository) GetArtists(w http.ResponseWriter, r *http.Request) {
	artists, err := m.DB.GetArtists()

	returnAsJSON(artists, w, err)
}

// SONGS
func (m *Repository) AddSong(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// get fields
	title := r.Form.Get("title")
	release_date := r.Form.Get("release_date")
	duration := r.Form.Get("duration")
	artist_id_string := r.Form.Get("artist_id")
	album_id_string := r.Form.Get("album_id")

	artist_id, err := strconv.Atoi(artist_id_string)
	if err != nil {

	}
	album_id, err := strconv.Atoi(album_id_string)
	if err != nil {

	}
	songToAdd := models.Song{
		Title:        title,
		Release_date: release_date,
		Duration:     duration,
		Album_id:     album_id,
		Artist_id:    artist_id,
	}

	addedSong, err := m.DB.AddSong(songToAdd)
	if err != nil {
		log.Println(err)
	}

	addedSong.Artist_name, err = m.DB.GetArtistName(artist_id)
	if err != nil {
		log.Println(err)
	}
	returnAsJSON(addedSong, w, err)
}

func (m *Repository) GetSongs(w http.ResponseWriter, r *http.Request) {
	songs, err := m.DB.GetSongs()

	returnAsJSON(songs, w, err)
}

func (m *Repository) GetSong(w http.ResponseWriter, r *http.Request) {
	x := chi.URLParam(r, "id")
	song, err := m.DB.GetSong(x)
	returnAsJSON(song, w, err)
}

func (m *Repository) AddSongToPlaylist(w http.ResponseWriter, r *http.Request) {

}

func (m *Repository) AddSongToAlbum(w http.ResponseWriter, r *http.Request) {

}

// i is the models.XYZ property
func returnAsJSON(i interface{}, w http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err)
	}
	j, _ := json.MarshalIndent(i, "", "   ")
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		log.Print(err)
	}
}
