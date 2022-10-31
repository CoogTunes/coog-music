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

func (m *Repository) GetArtists(w http.ResponseWriter, r *http.Request) {
	artists, err := m.DB.GetArtists()
	if err != nil {
		log.Println(err)
	}
	j, _ := json.MarshalIndent(artists, "", "   ")
	log.Println(string(j))
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		log.Print(err)
	}
}

func (m *Repository) GetUser(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	user, err := m.DB.GetUser(id)
	if err != nil {
		log.Println(err)
	}
	j, _ := json.MarshalIndent(user, "", "   ")
	// log.Println(string(j))
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		log.Print(err)
	}
}

func (m *Repository) AddSong(w http.ResponseWriter, r *http.Request) {

}
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
	newUser := models.Users{"-1", username, password, first_name, last_name, gender, boolAdmin}

	m.DB.AddUser(newUser)
}

func (m *Repository) AddArtist(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// get fields
	name := r.Form.Get("name")
	artist_id := r.Form.Get("artist_id")
	location := r.Form.Get("location")
	join_date := r.Form.Get("join_date")
	stringAdmin := r.Form.Get("admin")
	var songs []int
	// make sure admin is a boolean
	boolAdmin, err := strconv.ParseBool(stringAdmin)
	if err != nil {
		log.Println(err)
	}

	int_artist_id, err := strconv.Atoi(artist_id)
	if err != nil {
		log.Println(err)
	}
	// joindate and songs[] should be empty to start.
	artistToAdd := models.Artist{name, int_artist_id, location, join_date, songs, boolAdmin}

	addedArtist, err := m.DB.AddArtist(artistToAdd)
	if err != nil {
		log.Println(err)
	}

	j, _ := json.MarshalIndent(addedArtist, "", "   ")
	log.Println(string(j))
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		log.Print(err)
	}
}

func (m *Repository) AddSongToPlaylist(w http.ResponseWriter, r *http.Request) {

}
func (m *Repository) GetSong(w http.ResponseWriter, r *http.Request) {

	x := chi.URLParam(r, "id")
	song, err := m.DB.GetSong(x)
	if err != nil {
		log.Println(err)
	}
	j, _ := json.MarshalIndent(song, "", "   ")
	// log.Println(string(j))
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		log.Print(err)
	}
}
func (m *Repository) AddSongToAlbum(w http.ResponseWriter, r *http.Request) {

}
