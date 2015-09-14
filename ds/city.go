package ds

import (
	"github.com/boltdb/bolt"
	"strings"
)

type City struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CountryCode string `json:"country_code"`
	Population  string `json:"population"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Timezone    string `json:"timezone"`
}

type Cities struct {
	Cities []*City `json:"cities,omitempty"`
}

func cityFromString(id string, cityString string) *City {
	cityData := strings.Split(cityString, "\t")

	return &City{
		Id:          id,
		Name:        cityData[0],
		CountryCode: cityData[1],
		Population:  cityData[2],
		Latitude:    cityData[3],
		Longitude:   cityData[4],
		Timezone:    cityData[5],
	}
}

func FindCity(db *bolt.DB, id string) (*City, error) {
	var city *City = nil

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(CitiesBucketName)
		val := bucket.Get([]byte(id))

		if val != nil {
			city = cityFromString(id, string(val))
		}
		return nil
	})

	return city, err
}

func SearchCitiesByCityName(
	db *bolt.DB, locales []string, query string, limit int,
) (*Cities, error) {
	var cities Cities

	cityNames, err := SearchCityNames(db, locales, query, limit)

	var city *City
	for _, cityName := range *cityNames {
		city, err = FindCity(db, cityName.CityId)
		city.Name = cityName.Name
		cities.Cities = append(cities.Cities, city)
	}

	return &cities, err
}