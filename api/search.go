package api

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Rosya-edwica/go-trudvsem/models"
	"github.com/tidwall/gjson"
)

var Cities []models.City

func FindVacanciesByPosition(position models.Profession, cities []models.City) (positionVacancies []models.Vacancy){
	Cities = cities
	var pageNum = "0"
	
	for {
		params := url.Values{
			"text": {position.Name},
			"offset": {pageNum},
		}
		positionUrl := "http://opendata.trudvsem.ru/api/v1/vacancies?" + params.Encode()
		vacancies := ScrapePage(positionUrl, position)
		if len(vacancies) == 0 {
			break
		}
		num, _ := strconv.Atoi(pageNum)
		pageNum = strconv.Itoa(num + 1)
		positionVacancies = append(positionVacancies, vacancies...)
	}
	return

}


func ScrapePage(url string, position models.Profession) (vacancies []models.Vacancy) {
	json, err := GetJson(url)
	checkErr(err)

	for _, item := range gjson.Get(json, "results.vacancies").Array() {
		var vacancy models.Vacancy
		vacancy.Id = item.Get("vacancy.id").String()
		vacancy.Title = item.Get("vacancy.job-name").String()
		vacancy.ProfessionId = position.Id
		vacancy.SalaryFrom = item.Get("vacancy.salary_min").Float()
		vacancy.SalaryTo = item.Get("vacancy.salary_max").Float()
		vacancy.Url = item.Get("vacancy.vac_url").String()
		vacancy.ProfAreas = item.Get("vacancy.category.specialisation").String()
		vacancy.DateUpdate = item.Get("vacancy.creation-date").String()
		city := item.Get("vacancy.addresses.address.0.location").String()
		vacancy.CityId = parseCity(city)
		if vacancy.CityId == 0 {
			fmt.Println("Не смогли найти город в нашей БД", city)
		}

		vacancies = append(vacancies, vacancy)
	}

	return
}

func parseCity(cityName string) (cityId int) {
	re := regexp.MustCompile(`г .*?,|г .*? `)
	reSub := regexp.MustCompile(`г |г\.|,`)

	city := re.FindString(cityName + " ")
	if len(city) == 0 {
		return
	}
	
	city = reSub.ReplaceAllString(city, "")
	city = strings.TrimSpace(city)
	for _, item := range Cities {
		if strings.ToLower(item.Name) == strings.ToLower(city) {
			return item.EDWICA_ID
		}
	}
	
	return
}