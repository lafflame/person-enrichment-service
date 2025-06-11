package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

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
	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Инициализация БД
	db, err = initDB()
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatalf("Не удалось проверить соединение с БД: %v", err)
	}

	fmt.Println("Успешное подключение к PostgreSQL!")

	getName()

	var name string
	for {
		fmt.Println("\nВведите имя (или 'exit' для выхода):")
		fmt.Scan(&name)

		if name == "exit" {
			break
		}

		age := agify(name)
		if age != 0 {
			gender := genderize(name)
			nation := nationalize(name)
			addingDb(name, age, gender, nation)
		}
	}
}

func initDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	password := os.Getenv("DB_PASSWORD")

	// Проверяем, что все обязательные переменные установлены
	if host == "" || port == "" || user == "" || dbname == "" {
		return nil, fmt.Errorf("не все обязательные переменные окружения установлены")
	}

	// Преобразуем порт в число
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("неверный формат порта: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		host, p, user, dbname, sslmode, password)

	return sql.Open("postgres", connStr)
}

func addingDb(name string, age int, gender string, nation string) {
	initDB()
	defer db.Close()

	add, err := db.Exec("insert into users (name, age, gender, nation) values ($1, $2, $3, $4)", name, age, gender, nation)
	if err != nil {
		fmt.Println("Не удалось записать данные в таблицу БД, попробуйте заново")
	}
	fmt.Println(add)
}

// TODO
func getName() {
	// Подключение к базе данных
	initDB()
	defer db.Close()

	// Получение имени от пользователя
	var name string
	fmt.Println("Введите имя для поиска по базе данных:")
	fmt.Scan(&name)

	// Выполнение запроса и получение одной записи
	var id int
	var foundName string
	var age int
	var gender string

	err := db.QueryRow("SELECT id, name, age, gender FROM users WHERE name = $1", name).Scan(
		&id, &foundName, &age, &gender)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Пользователь с именем '%s' не найден\n", name)
		} else {
			fmt.Println("Не удалось подключиться к БД, попробуйте заново")
		}
		return
	}

	// Вывод найденного пользователя
	fmt.Printf("Найден пользователь:\nID: %d\nИмя: %s\nВозраст: %d\nПол: %s\n",
		id, foundName, age, gender)
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

	if result.Age == 0 {
		fmt.Println("Данных об этом имени недостаточно")
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
