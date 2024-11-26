package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"music/tools"
	"net/url"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type SongData struct {
	Song        string    `json:"song"`
	Group       string    `json:"group"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

var config *tools.Config = tools.GetConfig()

// Открывает соединение с БД
func OpenConnection(config *tools.Config) (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.DBName)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println("Failed to connect to the database: ", err)
		return nil, err
	}
	log.Println("Database connection opened")
	return db, nil
}

// Проверяет существование песни в БД
func Exists(song, group string) (int, error) {
	db, err := OpenConnection(config)
	if err != nil {
		return -1, err
	}
	defer db.Close()
	defer log.Println("Database connection closed")

	userID := -1
	statement := fmt.Sprintf(`SELECT s.song_id FROM "Song" s JOIN "Group" g ON s.group_id = g.group_id WHERE s.name = '%s' AND g.name = '%s'`, song, group)
	rows, err := db.Query(statement)
	if err != nil {
		log.Println("Failed to execute SELECT query: ", err)
		return -1, err
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Println("Failed to scan from sql.Rows: ", err)
			return -1, err
		}
		userID = id
	}
	defer rows.Close()
	return userID, nil
}

// Получает данные о песни по id
func GetSong(id int) (SongData, error) {
	data := SongData{}
	db, err := OpenConnection(config)
	if err != nil {
		return data, err
	}
	defer db.Close()
	defer log.Println("Database connection closed")

	statement := fmt.Sprintf(`SELECT s.name, g.name, "release_date", "text", "link" FROM "Song" s JOIN "Group" g on s.group_id = g.group_id WHERE song_id = %d`, id)
	rows, err := db.Query(statement)
	if err != nil {
		log.Println("Failed to execute SELECT query: ", err)
		return data, err
	}

	for rows.Next() {
		var dateString string
		err = rows.Scan(&data.Song, &data.Group, &dateString, &data.Text, &data.Link)
		if err != nil {
			log.Println("Failed to scan from sql.Rows: ", err)
			return data, err
		}
		data.ReleaseDate, err = time.Parse("2006-01-02T15:04:05Z07:00", dateString)
		if err != nil {
			log.Println("Failed to parse time: ", err)
			return data, err
		}

	}
	defer rows.Close()
	return data, nil
}

// Конструирует запрос на основе фильтра
func BuildListQuery(params url.Values) string {
	query := `SELECT s.name song, g.name "group", "release_date", "text", "link" FROM "Song" s JOIN "Group" g on s.group_id = g.group_id `
	emptyParams := true

	for param, list := range params {
		if len(list) > 0 && param != "page" {
			query += " WHERE "
			break
		}
	}

	for param, list := range params {
		if param != "page" && len(list) != 0 {
			if param == "releasedate" {
				param = "release_date"
			}
			emptyParams = false
			if param == "song" {
				query += `s.name IN ('`
			} else if param == "group" {
				query += `g.name IN ('`
			} else {
				query += `'` + param + `' IN ('`
			}
			for _, value := range list {
				if param == "release_date" {
					parsedDate, _ := time.Parse("02.01.2006", value)
					value = parsedDate.Format("2006-01-02")
				}
				query += value + "', '"
			}
			query = query[:len(query)-3]
			query += ") AND "
		}
	}

	if !emptyParams {
		query = query[:len(query)-4]
	}

	if len(params["page"]) != 0 {
		page, _ := strconv.Atoi(params["page"][0])
		query += fmt.Sprintf("LIMIT 10 OFFSET %d", (page-1)*10)
	}

	return query
}

// Добавляет новую песню
func AddSong(data SongData) error {
	db, err := OpenConnection(config)
	if err != nil {
		return err
	}
	defer db.Close()
	defer log.Println("Database connection closed")

	userID, err := Exists(data.Song, data.Group)
	if err != nil {
		return err
	}
	if userID != -1 {
		log.Printf("Attempt to add an existing song: '%s' by '%s'\n", data.Song, data.Group)
		err = errors.New("song already exists")
		return err
	}

	statement1 := `
		INSERT INTO "Group" (name)
		SELECT CAST($1 AS VARCHAR)
		WHERE NOT EXISTS (
    		SELECT 1 
    		FROM "Group" 
    		WHERE name = $1
		);`

	_, err = db.Exec(statement1, data.Group)
	if err != nil {
		log.Println("Failed to execute INSERT query 1: ", err)
		return err
	}

	statement2 := `INSERT INTO "Song" ("name", "release_date", "text", "link", "group_id")
		VALUES	($2, $3, $4, $5, 
			(
			SELECT group_id
			FROM "Group"
			WHERE name = $1
			)
				);`

	_, err = db.Exec(statement2, data.Group, data.Song, data.ReleaseDate, data.Text, data.Link)
	if err != nil {
		log.Println("Failed to execute INSERT query 2: ", err)
		return err
	}
	log.Printf("Song '%s' by '%s' added successfully\n", data.Song, data.Group)
	return nil
}

// Удаляет песню
func DeleteSong(song, group string) error {
	db, err := OpenConnection(config)
	if err != nil {
		return err
	}
	defer db.Close()
	defer log.Println("Database connection closed")

	id, err := Exists(song, group)
	if err != nil {
		return err
	}
	if id == -1 {
		log.Printf("Attempt to delete a non-existent song: '%s' by '%s'\n", song, group)
		err = errors.New("song does not exist")
		return err
	}

	statement := `delete from "Song" where song_id = $1`
	_, err = db.Exec(statement, id)
	if err != nil {
		log.Println("Failed to execute DELETE query: ", err)
		return err
	}

	log.Printf("Song '%s' by '%s' deleted successfully\n", song, group)
	return nil
}

// Обновляет информацию о песне
func UpdateSong(data SongData) error {
	db, err := OpenConnection(config)
	if err != nil {
		return err
	}
	defer db.Close()
	defer log.Println("Database connection closed")

	id, err := Exists(data.Song, data.Group)
	if err != nil {
		return err
	}
	if id == -1 {
		log.Printf("Attempt to update a non-existent song: '%s' by '%s'\n", data.Song, data.Group)
		err = errors.New("song does not exist")
		return err
	}

	oldData, err := GetSong(id)
	if err != nil {
		return err
	}

	if data.ReleaseDate.IsZero() {
		data.ReleaseDate = oldData.ReleaseDate
	}
	if data.Text == "" {
		data.Text = oldData.Text
	}
	if data.Link == "" {
		data.Link = oldData.Link
	}

	statement := `update "Song" set "release_date" = $1, "text" = $2, "link" = $3 where "song_id" = $4`
	_, err = db.Exec(statement, data.ReleaseDate, data.Text, data.Link, id)
	if err != nil {
		log.Println("Failed to execute UPDATE query: ", err)
		return err
	}

	log.Printf("Song '%s' by '%s' updated successfully\n", data.Song, data.Group)
	return nil
}

// Получает текст песни
func GetText(song, group string) (string, error) {
	id, err := Exists(song, group)
	if err != nil {
		return "", err
	}
	if id == -1 {
		log.Printf("Attempt to get text of non-existent song: '%s' by '%s'\n", song, group)
		err = errors.New("song does not exist")
		return "", err
	}

	data, err := GetSong(id)
	if err != nil {
		return "", err
	}

	text := data.Text

	log.Printf("Got text of '%s' by '%s' successfully\n", song, group)
	return text, nil
}

// Получает список песен
func ListSongs(params url.Values) ([]SongData, error) {
	data := []SongData{}
	db, err := OpenConnection(config)
	if err != nil {
		return data, err
	}
	defer db.Close()
	defer log.Println("Database connection closed")

	statment := BuildListQuery(params)
	rows, err := db.Query(statment)
	if err != nil {
		log.Println("Failed to execute SELECT query ", err)
		return data, err
	}

	for rows.Next() {
		temp := SongData{}
		dateString := ""
		err = rows.Scan(&temp.Song, &temp.Group, &dateString, &temp.Text, &temp.Link)
		if err != nil {
			log.Println("Failed to scan sql.Rows: ", err)
			return data, err
		}
		temp.ReleaseDate, err = time.Parse("2006-01-02T15:04:05Z07:00", dateString)
		if err != nil {
			log.Println("Failed ti parse time: ", err)
			return data, err
		}
		data = append(data, temp)

	}

	log.Println("Got list of songs successfully")
	return data, nil
}
