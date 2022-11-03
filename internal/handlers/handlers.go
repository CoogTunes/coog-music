package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/DeLuci/coog-music/internal/forms"
	"github.com/DeLuci/coog-music/internal/render"
	"golang.org/x/crypto/bcrypt"

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

// LOGIN/SIGNUP

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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
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
	}
	users, err := m.DB.AddUser(registrationModel)
	if false {
		log.Println(users)
	}
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// HOME PAGE

func (m *Repository) GetHome(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "coogtunes.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// END OF HOME PAGE

// PROFILE PAGE
func (m *Repository) GetProfile(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "profilepage.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
	})
}

//  END OF PROFILE PAGE

// USERS
func (m *Repository) AddUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// get fields
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	first_name := r.Form.Get("first_name")
	last_name := r.Form.Get("last_name")
	stringAdmin := r.Form.Get("admin")

	admin, err := strconv.Atoi(stringAdmin)
	if err != nil {
		log.Println(err)
	}

	newUser := models.Users{Username: username, Password: password, First_name: first_name, Last_name: last_name, Admin_level: admin}

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
	var songs []models.Song

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

func (m *Repository) GetArtistsAndSongs(w http.ResponseWriter, r *http.Request) {
	artists, err := m.DB.GetArtistsAndSongs()
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

func (m *Repository) UpdateSongCount(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
	}
	var songWithId models.Song = models.Song{Song_id: id}
	song, err := m.DB.UpdateSongCount(songWithId)
	returnAsJSON(song, w, err)
}

func (m *Repository) AddSongToPlaylist(w http.ResponseWriter, r *http.Request) {

}

func (m *Repository) AddSongToAlbum(w http.ResponseWriter, r *http.Request) {

}

func (m *Repository) GetPlaylists(w http.ResponseWriter, r *http.Request) {
	playlists, err := m.DB.GetPlaylists()

	returnAsJSON(playlists, w, err)
}

func (m *Repository) GetAlbums(w http.ResponseWriter, r *http.Request) {
	albums, err := m.DB.GetAlbums()

	returnAsJSON(albums, w, err)
}

func (m *Repository) UpdateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// get fields
	userIdString := r.Form.Get("user_id")
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	first_name := r.Form.Get("first_name")
	last_name := r.Form.Get("last_name")
	stringAdmin := r.Form.Get("admin_level")

	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		log.Println(err)
	}

	admin, err := strconv.Atoi(stringAdmin)
	if err != nil {
		log.Println(err)
	}

	newUser := models.Users{User_id: userId, Username: username, Password: password, First_name: first_name, Last_name: last_name, Admin_level: admin}

	addedUser, err := m.DB.UpdateUser(newUser)

	returnAsJSON(addedUser, w, err)
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
