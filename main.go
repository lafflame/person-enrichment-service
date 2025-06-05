package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/lib/pq"
)

type AgifyResponse struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type GenderizeResponse struct {
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type NationalizeResponse struct {
	Name    string    `json:"name"`
	Country []Country `json:"country"`
}

type Country struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

func main() {
	var name string
	for {
		fmt.Println("\nВведите имя:")
		fmt.Scan(&name)
		age := agify(name)
		gender := genderize(name)
		nation := nationalize(name)
		addingDb(name, age, gender, nation)
	}

}

func agify(name string) int {
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		fmt.Println("Err")
		return 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return 0
	}

	var result AgifyResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
		return 0
	}

	fmt.Printf("Имя: %s\nВозраст: %d\n", result.Name, result.Age)

	return result.Age
}

func genderize(name string) string {
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		fmt.Println("Err")
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return ""
	}

	var result GenderizeResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
		return ""
	}

	fmt.Printf("Пол: %s (вероятность: %.2f%%)\n", result.Gender, result.Probability*100)

	return result.Gender
}

func nationalize(name string) string {
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		fmt.Println("Err")
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return ""
	}

	var result NationalizeResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
		return ""
	}

	fmt.Printf("Национальность: %s \n", result.Country[0].CountryID)

	return result.Country[0].CountryID
}

func addingDb(name string, age int, gender string, nation string) {
	connect := "host=localhost port=5432 user=postgres dbname=users sslmode=disable"
	db, err := sql.Open("postgres", connect)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	add, err := db.Exec("insert into users (name, age, gender, nation) values ($1, $2, $3, $4)", name, age, gender, nation)
	if err != nil {
		panic(err)
	}
	fmt.Println(add.RowsAffected())
}
