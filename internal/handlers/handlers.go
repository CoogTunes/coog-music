package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/DeLuci/coog-music/internal/forms"
	"github.com/DeLuci/coog-music/internal/models"
	"github.com/DeLuci/coog-music/internal/render"
	"golang.org/x/crypto/bcrypt"
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

func (m *Repository) GetLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	loginOrRegister := r.Form.Get("submit_button")

	if loginOrRegister == "register" {
		fmt.Println("Passes to function PostRegistration ")
		m.PostRegistration(w, r)
		return
	}
	//_ = m.App.Session.RenewToken(r.Context())

	email := r.Form.Get("email")
	pwd := r.Form.Get("password")

	err = m.DB.Authenticate(email, pwd)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) PostRegistration(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	pwd := []byte(r.Form.Get("password"))
	hashedPwd, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return
	}

	registrationModel := models.Users{
		First_name: r.Form.Get("first_name"),
		Last_name:  r.Form.Get("last_name"),
		Username:   r.Form.Get("email"),
		Password:   string(hashedPwd),
		Gender:     "Male",
	}
	err = m.DB.AddUser(registrationModel)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) AddSong(w http.ResponseWriter, r *http.Request) {

}
func (m *Repository) AddUser(w http.ResponseWriter, r *http.Request) {

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
