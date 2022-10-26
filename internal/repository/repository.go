package repository

import (
	"github.com/DeLuci/coog-music/internal/models"
)

type DatabaseRepo interface {
	GetArtists() ([]models.Artist, error)
}
