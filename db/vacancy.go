package db

import (
	"fmt"
	"strings"

	"github.com/Rosya-edwica/go-trudvsem/models"
)

func (d *Database) SaveVacancy(v models.Vacancy) (err error){
	if v.Title == "" {
		return
	}

	columns := buildPatternInsertValues(12)
	smt := fmt.Sprintf(`INSERT INTO h_vacancy (id, url, name, city_id, position_id, prof_areas, specs, experience, salary_from, salary_to, key_skills, vacancy_date, platform) VALUES %s`, columns)
	tx, _ := d.Connection.Begin()
	_, err = d.Connection.Exec(smt, v.Id, v.Url, v.Title, v.CityId, v.ProfessionId, v.ProfAreas, v.Specializations, v.Experience, v.SalaryFrom, v.SalaryTo, v.Skills, v.DateUpdate, "trudvsem")
	if err != nil {
		fmt.Printf("Ошибка: Вакансия %s не была добавлена в базу - %s\n", v.Id, err)
		return
	}
	err = tx.Commit()
	checkErr(err)
	fmt.Println("Успешно сохранили вакансию", v.Title)
	return nil
}


func buildPatternInsertValues(valuesCount int) (pattern string) {
	var items []string
	
	for i := 0; i < valuesCount; i++ {
		items = append(items, "?")
	}
	pattern = strings.Join(items, ",")
	return fmt.Sprintf("(%s)", pattern)
}

func (d *Database) SaveManyVacancies(vacancies []models.Vacancy) {
	groups := groupVacancies(vacancies)
	for _, group := range groups {
		d.SaveVacancies(group)
	}
}

func (d *Database) SaveVacancies(vacancies []models.Vacancy) {
	if len(vacancies) == 0 { return }
	query, vals := createQueryForMultipleInsertVacanciesMYSQL(vacancies)
	tx, _ := d.Connection.Begin()
	_, err := d.Connection.Exec(query, vals...)
	checkErr(err)
	tx.Commit()
	fmt.Printf("Успешко сохранили %d вакансий \n", len(vacancies))
}


func groupVacancies(vacancies []models.Vacancy) (groups [][]models.Vacancy) {
	LIMIT := 2000
	for i := 0; i < len(vacancies); i += LIMIT {
		group := vacancies[i:]
		if len(group) >= 2000 {
			group = group[:LIMIT]
		}
		groups = append(groups, group)
	}
	return
}

func createQueryForMultipleInsertVacanciesMYSQL(vacancies []models.Vacancy) (query string, valArgs []interface{}) {
	query = "INSERT IGNORE INTO h_vacancy (id, name, url, city_id, position_id, prof_areas, specs, experience, salary_from, salary_to, key_skills, vacancy_date, platform) VALUES "
	for _, v := range vacancies {
		query += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
		valArgs = append(valArgs,  v.Id, v.Title, v.Url, v.CityId, v.ProfessionId, v.ProfAreas, v.Specializations, v.Experience, v.SalaryFrom, v.SalaryTo, v.Skills, v.DateUpdate, "trudvsem")
	}
	query = query[0:len(query)-1]
	query = strings.TrimSpace(query)
	return
}