package repository

import (
	"github.com/DeLuci/coog-music/internal/models"
)

type DatabaseRepo interface {
	GetArtists2() ([]models.Artist, error)
}
