package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	Count   int       `json:"count"`
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
		fmt.Scan(&name)
		agify(name)
		genderize(name)
		nationalize(name)
	}
}

func agify(name string) {
	resp, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		fmt.Println("Err")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return
	}

	var result AgifyResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
		return
	}

	fmt.Printf("Имя: %s\nВозраст: %d\n", result.Name, result.Age)
}

func genderize(name string) {
	resp, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		fmt.Println("Err")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return
	}

	var result GenderizeResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
		return
	}

	fmt.Printf("Пол: %s (вероятность: %.2f%%)\n", result.Gender, result.Probability*100)
}

func nationalize(name string) {
	resp, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		fmt.Println("Err")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return
	}

	var result NationalizeResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
		return
	}

	fmt.Printf("Национальность для %s:\n", result.Name)
	for i, country := range result.Country {
		fmt.Printf("%d. Страна: %s, Вероятность: %.2f%%\n",
			i+1,
			country.CountryID,
			country.Probability*100)
	}
}
