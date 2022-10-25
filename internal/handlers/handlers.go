package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DeLuci/coog-music/internal/config"
	"github.com/DeLuci/coog-music/internal/driver"

	// "github.com/DeLuci/coog-music/internal/models"
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

	// m.DB.GetArtists()
	// json.NewEncoder(w).Encode()

	// productParam := chi.URLParam(r, "product")
	//decodedValue := url.QueryEscape(productParam)
	//splitProduct := strings.Split(productParam, " ")
	artists, err := m.DB.GetArtists()
	if err != nil {
		log.Println("START")
		log.Println(err)
		log.Println("END")
	}
	j, _ := json.MarshalIndent(artists, "", "   ")
	log.Println(string(j))
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		log.Print(err)
	}
}
