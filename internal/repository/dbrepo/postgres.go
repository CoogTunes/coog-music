package dbrepo

import (
	"github.com/DeLuci/coog-music/internal/models"
)

func (m *postgresDBRepo) GetArtists() ([]models.Artist, error) {

	var artists []models.Artist
	var artist models.Artist

	artist.Name = "me my name"

	// artists = append(artists, {})
	// query := "select * from artists"

	// rows, err := m.DB.Query(query)
	// if err != nil {
	// 	return nil, err
	// }

	// defer func(rows *sql.Rows) {
	// 	err := rows.Close()
	// 	if err != nil {

	// 	}
	// }(rows)

	// for rows.Next() {
	// 	var artist models.Artist

	// 	rows.Scan(&artist.Name, &artist.Artist_id, &artist.Location, &artist.Join_date, &artist.Songs, &artist.Admin)

	// 	if err != nil {
	// 		return nil, err
	// 	}
	artists = append(artists, artist)
	// }
	return artists, nil
}
