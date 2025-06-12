# Программа для определения возраста, пола и национальности по имени

Это приложение на Go, которое использует внешние API для определения вероятного возраста, пола и национальности по имени, а также сохраняет результаты в PostgreSQL базе данных.

## Функциональность

- Получение предполагаемого возраста по имени (через API Agify)
- Определение вероятного пола по имени (через API Genderize)
- Определение вероятной национальности по имени (через API Nationalize)
- Сохранение результатов в PostgreSQL базе данных
- Поиск ранее сохраненных данных по имени

## Требования

- Go 1.16+
- PostgreSQL
- Доступ к интернету (для работы с внешними API)

## Установка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/lafflame/person-enrichment-service.git
   cd person-enrichment-service
   ```

2. Установите зависимости:
   ```bash
   go mod download
   ```

3. Создайте файл `.env` в корне проекта с вашими настройками PostgreSQL:
   ```
   DB_HOST=your_host
   DB_PORT=your_port
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=your_database
   DB_SSLMODE=disable
   ```

4. Создайте таблицу в вашей PostgreSQL базе данных:
   ```sql
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       name VARCHAR(100) NOT NULL,
       age INTEGER,
       gender VARCHAR(20),
       nation VARCHAR(50)
   );
   ```

## Использование

1. Запустите программу:
   ```bash
   go run main.go
   ```

2. Следуйте инструкциям в терминале:
   - Введите имя для анализа
   - Программа покажет предполагаемые возраст, пол и национальность
   - Результаты будут сохранены в базу данных
   - Можно искать ранее сохраненные имена в базе

3. Для выхода введите `exit`

## API Использованные в проекте

- [Agify.io](https://agify.io/) - для определения возраста по имени
- [Genderize.io](https://genderize.io/) - для определения пола по имени
- [Nationalize.io](https://nationalize.io/) - для определения национальности по имени

## Структура кода

- `main.go` - основной файл с точкой входа
- Функции:
  - `initDB()` - инициализация подключения к PostgreSQL
  - `addingDb()` - добавление данных в базу
  - `getName()` - поиск имени в базе
  - `agify()` - запрос к API возраста
  - `genderize()` - запрос к API пола
  - `nationalize()` - запрос к API национальности
