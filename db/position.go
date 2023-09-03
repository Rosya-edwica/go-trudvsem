package db

import (
	"database/sql"
	"strings"

	"github.com/Rosya-edwica/go-trudvsem/models"
)

func (d *Database) GetPositions() (positions []models.Profession) {
	query := "SELECT id, name, other_names FROM position"
	rows, err := d.Connection.Query(query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var (
			name string
			other sql.NullString
			id int
		)
		err = rows.Scan(&id, &name, &other)
		checkErr(err)

		prof := models.Profession{
			Id:         id,
			Name:       name,
			OtherNames: strings.Split(other.String, "|"),
		}
		positions = append(positions, prof)

	}
	return
}