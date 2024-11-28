package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"music/internal/database"
	"music/tools"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SongData struct {
	Song        string `json:"song"`
	Group       string `json:"group"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// Получает список песен
func GetSongs(params url.Values) ([]database.SongData, []string, error) {
	songs := []database.SongData{}

	expectedParams := map[string]bool{
		"song":        true,
		"group":       true,
		"releasedate": true,
		"text":        true,
		"link":        true,
		"page":        true,
		"onpage":      true,
	}

	var unexpectedParams []string

	for param := range params {
		params[param] = strings.Split(params[param][0], ",")
	}

	// Проверка на лишние параметры
	for param := range params {
		if _, ok := expectedParams[param]; !ok {
			unexpectedParams = append(unexpectedParams, param)
		}
	}

	// Если нашли лишние параметры
	if len(unexpectedParams) != 0 {
		tools.Logger.Info(fmt.Sprintf("Unexpected parameters passed: %s", strings.Join(unexpectedParams, ", ")))
		err := errors.New("unexpected params")
		return songs, unexpectedParams, err
	}

	// Валидация параметра page
	if len(params["page"]) > 1 {
		tools.Logger.Info(fmt.Sprintf("Invalid 'page' format passed: %s", params["page"]))
		err := errors.New("'page' requires only 1 value")
		return songs, unexpectedParams, err
	} else if len(params["page"]) != 0 {
		page, err := strconv.Atoi(params["page"][0])
		if err != nil || page < 1 {
			tools.Logger.Info(fmt.Sprintf("Invalid 'page' format passed: %s", params["page"][0]))
			err := errors.New("page is not a number")
			return songs, unexpectedParams, err
		}
	}

	// Валидация параметра onpage
	if len(params["onpage"]) > 1 {
		tools.Logger.Info(fmt.Sprintf("Invalid 'onpage' format passed: %s", params["onpage"]))
		err := errors.New("'onpage' requires only 1 value")
		return songs, unexpectedParams, err
	} else if len(params["onpage"]) != 0 {
		onpage, err := strconv.Atoi(params["onpage"][0])
		if err != nil || onpage < 1 {
			tools.Logger.Info(fmt.Sprintf("Invalid 'onpage' format passed: %s", params["onpage"][0]))
			err := errors.New("onpage is not a number")
			return songs, unexpectedParams, err
		}
	}

	// Валидация формата даты
	if date, ok := params["releasedate"]; ok {
		_, err := time.Parse("2.1.2006", date[0])
		if err != nil {
			tools.Logger.Info(fmt.Sprintf("Invalid date format passed: %s", date[0]))
			err := errors.New("incorrect date format")
			return songs, unexpectedParams, err
		}
	}

	songs, err := database.ListSongs(params)
	if err != nil {
		return songs, unexpectedParams, err
	}
	return songs, unexpectedParams, nil

}

// Получает текст песни
func GetText(params url.Values) (string, []string, error) {
	text := ""

	expectedParams := map[string]bool{
		"song":  true,
		"group": true,
		"verse": true,
	}

	requiredParams := map[string]bool{
		"song":  true,
		"group": true,
	}

	var unexpectedParams []string

	for param := range params {
		params[param] = strings.Split(params[param][0], ",")
	}

	// Проверка на лишние параметры
	for param := range params {
		if _, ok := expectedParams[param]; !ok {
			unexpectedParams = append(unexpectedParams, param)
		}
	}

	// Если нашли лишние параметры
	if len(unexpectedParams) != 0 {
		tools.Logger.Info(fmt.Sprintf("Unexpected parameters passed: %s", strings.Join(unexpectedParams, ", ")))
		err := errors.New("unexpected params")
		return text, unexpectedParams, err
	}

	// Валидация пераметров song и group
	for param := range requiredParams {
		if _, ok := params[param]; !ok {
			tools.Logger.Info(fmt.Sprintf("Required parameter '%s' was not passed\n", param))
			errorMessage := fmt.Sprintf("'%s' parameter is required", param)
			err := errors.New(errorMessage)
			return text, unexpectedParams, err
		} else if len(params[param]) != 1 {
			tools.Logger.Info(fmt.Sprintf("To many '%s' parameters was passed\n", param))
			errorMessage := fmt.Sprintf("'%s' requires only 1 value", param)
			err := errors.New(errorMessage)
			return text, unexpectedParams, err
		}
	}

	// Валидация параметра verse
	if len(params["verse"]) > 1 {
		tools.Logger.Info("To many 'verse' parameters was passed")
		err := errors.New("'verse' requires only 1 value")
		return text, unexpectedParams, err
	} else if len(params["verse"]) != 0 {
		verse, err := strconv.Atoi(params["verse"][0])
		if err != nil || verse < 1 {
			tools.Logger.Info("Invalid 'verse' format was passed")
			err := errors.New("'verse' requires a positive number")
			return text, unexpectedParams, err
		}
	}

	// Получаем текст песни
	text, err := database.GetText(params["song"][0], params["group"][0])
	if err != nil {
		if err.Error() == "song does not exist" {
			return text, unexpectedParams, err
		} else {
			err = errors.New("failed to get songs text")
			return text, unexpectedParams, err
		}
	}

	// Разбиваем по куплетам
	if len(params["verse"]) == 0 {
		return text, unexpectedParams, nil
	} else {
		verses := strings.Split(text, "\n\n")
		verse, _ := strconv.Atoi(params["verse"][0])
		if verse > len(verses) {
			return "", unexpectedParams, nil
		}

		return verses[verse-1], unexpectedParams, nil
	}

}

// Удаляет песню
func DeleteSong(params url.Values) ([]string, error) {
	requiredParams := map[string]bool{
		"song":  true,
		"group": true,
	}

	var unexpectedParams []string

	for param := range params {
		params[param] = strings.Split(params[param][0], ",")
	}

	// Проверка на лишние параметры
	for param := range params {
		if _, ok := requiredParams[param]; !ok {
			unexpectedParams = append(unexpectedParams, param)
		}
	}

	// Если нашли лишние параметры
	if len(unexpectedParams) != 0 {
		tools.Logger.Info(fmt.Sprintf("Unexpected parameters passed: %s", strings.Join(unexpectedParams, ", ")))
		err := errors.New("unexpected params")
		return unexpectedParams, err
	}

	// Валидация пераметров
	for param := range requiredParams {
		if _, ok := params[param]; !ok {
			tools.Logger.Info(fmt.Sprintf("Required parameter '%s' was not passed\n", param))
			errorMessage := fmt.Sprintf("'%s' parameter is required", param)
			err := errors.New(errorMessage)
			return unexpectedParams, err
		} else if len(params[param]) != 1 {
			tools.Logger.Info(fmt.Sprintf("To many '%s' parameters was passed\n", param))
			errorMessage := fmt.Sprintf("'%s' requires only 1 value", param)
			err := errors.New(errorMessage)
			return unexpectedParams, err
		}
	}

	// Удаление песни
	err := database.DeleteSong(params["song"][0], params["group"][0])
	if err != nil {
		if err.Error() == "song does not exist" {
			return unexpectedParams, err
		} else {
			err = errors.New("failed to delete song")
			return unexpectedParams, err
		}
	} else {
		return unexpectedParams, nil
	}
}

// Обновляет информацию о песни
func UpdateSong(params map[string]string) ([]string, error) {
	expectedParams := map[string]bool{
		"song":        true,
		"group":       true,
		"releasedate": true,
		"text":        true,
		"link":        true,
		"page":        true,
	}

	requiredParams := map[string]bool{
		"song":  true,
		"group": true,
	}

	unexpectedParams := []string{}

	// Проверка на лтшние параметры
	for param := range params {
		if _, ok := expectedParams[param]; !ok {
			unexpectedParams = append(unexpectedParams, param)
		}
	}

	// Если нашли лишние параметры
	if len(unexpectedParams) != 0 {
		tools.Logger.Info(fmt.Sprintf("Unexpected parameters passed: %s", strings.Join(unexpectedParams, ", ")))
		err := errors.New("unexpected params")
		return unexpectedParams, err
	}

	// Проверка на обязательные параметры
	for param := range requiredParams {
		if _, ok := params[param]; !ok {
			tools.Logger.Info(fmt.Sprintf("Required parameter '%s' was not passed\n", param))
			errorMessage := fmt.Sprintf("'%s' parameter is required", param)
			err := errors.New(errorMessage)
			return unexpectedParams, err

		}
	}

	// Заполнение структуры SongData
	var songData database.SongData
	for param, value := range params {
		switch param {
		case "song":
			songData.Song = value
		case "group":
			songData.Group = value
		case "releasedate":
			releaseDate, err := time.Parse("2.1.2006", value)
			if err != nil {
				tools.Logger.Info(fmt.Sprintf("Invalid date format: %s", value))
				err = errors.New("incorrect date format")
				return unexpectedParams, err
			}
			songData.ReleaseDate = releaseDate
		case "text":
			songData.Text = value
		case "link":
			songData.Link = value
		}
	}

	// Обновление данных о песни
	err := database.UpdateSong(songData)
	if err != nil {
		if err.Error() == "song does not exist" {
			return unexpectedParams, err
		} else {
			err = errors.New("failed to update song")
			return unexpectedParams, err
		}
	}

	return unexpectedParams, nil
}

// Добавляет новую песню
func AddSong(params url.Values) ([]string, error) {
	requiredParams := map[string]bool{
		"song":  true,
		"group": true,
	}

	unexpectedParams := []string{}

	// Проверка на лтшние параметры
	for param := range params {
		if _, ok := requiredParams[param]; !ok {
			unexpectedParams = append(unexpectedParams, param)
		}
	}

	// Если нашли лишние параметры
	if len(unexpectedParams) != 0 {
		tools.Logger.Info(fmt.Sprintf("Unexpected parameters passed: %s", strings.Join(unexpectedParams, ", ")))
		err := errors.New("unexpected params")
		return unexpectedParams, err
	}

	// Проверка на обязательные параметры
	for param := range requiredParams {
		if _, ok := params[param]; !ok {
			tools.Logger.Info(fmt.Sprintf("Required parameter '%s' was not passed\n", param))
			errorMessage := fmt.Sprintf("'%s' parameter is required", param)
			err := errors.New(errorMessage)
			return unexpectedParams, err

		}
	}

	// Получаем информацию о песни
	data, err := getSongInfo(params["song"][0], params["group"][0])
	if err != nil {
		tools.Logger.Error(fmt.Sprintf("Failed to get song info: '%s' by '%s'\n", params["song"], params["group"]), err)
		err = errors.New("failed to get song info")
		return unexpectedParams, err
	}
	songData := stringToDate(data)
	err = database.AddSong(songData)
	if err != nil {
		if err.Error() == "song already exists" {
			return unexpectedParams, err
		}
		err = errors.New("failed to add song")
		return unexpectedParams, err
	}

	return unexpectedParams, nil

}

// Делает запрос к music info API
func getSongInfo(song, group string) (SongData, error) {
	songData := SongData{}

	baseURL := tools.GetConfig().MusicInfoAddr
	if baseURL == "mock" {
		baseURL = "http://0.0.0.0:8181/info"
	}
	params := url.Values{}
	params.Add("group", group)
	params.Add("song", song)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		tools.Logger.Error("Failed to connect to song info API: ", err)
		return songData, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tools.Logger.Error("Got bad response from song info API: code ", errors.New(string(resp.StatusCode)))
		return songData, errors.New("failed to get song info")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tools.Logger.Error("Failed to read body from song info API response: ", err)
		return songData, err
	}
	err = json.Unmarshal(body, &songData)
	if err != nil {
		tools.Logger.Error("Failed to unmarshal body from song info API response: ", err)
		return songData, err
	}

	songData.Song = song
	songData.Group = group

	tools.Logger.Info(fmt.Sprintf("Got song info successfully: '%s' by '%s'\n", song, group))
	return songData, nil
}

// Вспомогательная функция
func DateToString(songs []database.SongData) []SongData {
	var result []SongData
	for _, song := range songs {
		temp := SongData{}
		temp.Song = song.Song
		temp.Group = song.Group
		temp.ReleaseDate = song.ReleaseDate.Format("02.01.2006")
		temp.Text = song.Text
		temp.Link = song.Link
		result = append(result, temp)
	}
	return result
}

// Вспомогательная функция
func stringToDate(data SongData) database.SongData {
	var result database.SongData
	result.Song = data.Song
	result.Group = data.Group
	result.ReleaseDate, _ = time.Parse("02.01.2006", data.ReleaseDate)
	result.Text = data.Text
	result.Link = data.Link

	return result
}
