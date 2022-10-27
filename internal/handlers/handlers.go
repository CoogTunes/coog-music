package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DeLuci/coog-music/internal/config"
	"github.com/DeLuci/coog-music/internal/driver"
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

func (m *Repository) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := m.DB.GetUsers()
	if err != nil {
		log.Println(err)
	}
	j, _ := json.MarshalIndent(users, "", "   ")
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
	// temp proof of concept for post requests
	log.Println("in add user post")
	r.ParseForm()

	pwd := r.Form.Get("password")
	email := r.Form.Get("email")

	log.Println("pwd, email", pwd, email)

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
