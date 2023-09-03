package db

import "github.com/Rosya-edwica/go-trudvsem/models"

func (d *Database) GetCities() (cities []models.City) {
	query := "SELECT id_edwica, name FROM h_city WHERE id_edwica != 0"

	rows, err := d.Connection.Query(query)
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		var name string
		var edwica_id int
		err = rows.Scan(&edwica_id, &name)
		checkErr(err)
		cities = append(cities, models.City{
			EDWICA_ID: edwica_id,
			Name:      name,
		})
	}
	return
}