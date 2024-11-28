package handlers

import (
	"encoding/json"
	"io"
	_ "music/api"
	"music/internal/services"
	"net/http"
	"strings"
	"time"
)

type SongData struct {
	Song        string    `json:"song"`
	Group       string    `json:"group"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

func SongsHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		query := request.URL.Query()
		songs, unexpectedParams, err := services.GetSongs(query)

		if err != nil {
			if err.Error() == "unexpected params" {
				errorMessage := "Unexpected parameters: " + strings.Join(unexpectedParams, ", ")
				http.Error(writer, errorMessage, http.StatusBadRequest)
				return
			} else if err.Error() == "page is not a number" {
				http.Error(writer, `"page" requires a positive number`, http.StatusBadRequest)
				return
			} else if err.Error() == "'page' requires only 1 value" {
				http.Error(writer, "'page' requires only 1 value", http.StatusBadRequest)
				return
			} else if err.Error() == "'onpage' requires only 1 value" {
				http.Error(writer, "'onpage' requires only 1 value", http.StatusBadRequest)
				return
			} else if err.Error() == "onpage is not a number" {
				http.Error(writer, `"onpage" requires a positive number`, http.StatusBadRequest)
				return
			} else if err.Error() == "incorrect date format" {
				http.Error(writer, "Invalid date format: "+query["releasedate"][0], http.StatusBadRequest)
				return
			} else {
				http.Error(writer, "Failed to get music list ", http.StatusInternalServerError)
				return
			}
		}

		writer.WriteHeader(200)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(services.DateToString(songs))
		return

	} else if request.Method == "DELETE" {
		query := request.URL.Query()
		unexpectedParams, err := services.DeleteSong(query)
		if err != nil {
			if err.Error() == "unexpected params" {
				errorMessage := "Unexpected parameters: " + strings.Join(unexpectedParams, ", ")
				http.Error(writer, errorMessage, http.StatusBadRequest)
				return
			} else if err.Error() == "failed to delete song" {
				http.Error(writer, "Failed to delete song", http.StatusInternalServerError)
				return
			} else if err.Error() == "song does not exist" {
				http.Error(writer, "Song does not exist", http.StatusNotFound)
				return
			} else {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}

		writer.WriteHeader(200)
		writer.Write([]byte("Song deleted"))
		return

	} else if request.Method == "PATCH" {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "Can't read request body", http.StatusBadRequest)
			return
		}
		defer request.Body.Close()

		var songUpdate map[string]string
		err = json.Unmarshal(body, &songUpdate)
		if err != nil {
			http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		unexpectedParams, err := services.UpdateSong(songUpdate)
		if err != nil {
			if err.Error() == "unexpected params" {
				errorMessage := "Unexpected parameters: " + strings.Join(unexpectedParams, ", ")
				http.Error(writer, errorMessage, http.StatusBadRequest)
				return
			} else if err.Error() == "failed to update song" {
				http.Error(writer, "Failed to update song data", http.StatusInternalServerError)
				return
			} else if err.Error() == "song does not exist" {
				http.Error(writer, "Song does not exist", http.StatusNotFound)
				return
			} else {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		writer.WriteHeader(200)
		writer.Write([]byte("Song data updated"))
		return

	} else if request.Method == "POST" {
		query := request.URL.Query()
		unexpectedParams, err := services.AddSong(query)
		if err != nil {
			if err.Error() == "unexpected params" {
				errorMessage := "Unexpected parameters: " + strings.Join(unexpectedParams, ", ")
				http.Error(writer, errorMessage, http.StatusBadRequest)
				return
			} else if err.Error() == "failed to add song" {
				http.Error(writer, "Failed to add song", http.StatusInternalServerError)
				return
			} else if err.Error() == "failed to get song info" {
				http.Error(writer, "Failed to get song info", http.StatusInternalServerError)
				return
			} else if err.Error() == "song already exists" {
				http.Error(writer, "Song already exists", http.StatusBadRequest)
				return
			} else {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}

		writer.WriteHeader(200)
		writer.Write([]byte("New song added"))
		return
	}
}

// Обработчик /text
func TextHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		query := request.URL.Query()
		text, unexpectedParams, err := services.GetText(query)

		if err != nil {
			if err.Error() == "unexpected params" {
				errorMessage := "Unexpected parameters: " + strings.Join(unexpectedParams, ", ")
				http.Error(writer, errorMessage, http.StatusBadRequest)
				return
			} else if err.Error() == "failed to get songs text" {
				http.Error(writer, "Failed to get songs text", http.StatusInternalServerError)
				return
			} else if err.Error() == "song does not exist" {
				http.Error(writer, "Song does not exist", http.StatusNotFound)
				return
			} else {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}

		writer.WriteHeader(200)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(text)
		return
	} else {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

// Обертки над SongHandler сделаны для генерации Swagger

// @Summary      Get a list of songs
// @Description  Get a list of songs based on filtering parameters.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        song        query    string  false  "Song name"
// @Param        group       query    string  false  "Group name"
// @Param        releasedate query   string  false  "Release date"
// @Param        text        query    string  false  "Song lyrics"
// @Param        link        query    string  false  "Video link"
// @Param        page        query    int     false  "Page number"
// @Param        onpage      query    int     false  "Items per page"
// @Success      200       {array}  SongData    "List of songs"
// @Failure      400        {string} string  "Bad request"
// @Failure      500        {string} string  "Internal server error"
// @Router       /songs [get]
func getSongHandler(w http.ResponseWriter, r *http.Request) {
	SongsHandler(w, r)
}

// @Summary      Add a new song
// @Description  Add a new song to the database.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        song        query    string  false  "Song name"
// @Param        group       query    string  false  "Group name"
// @Success      200   {string} string  "Song added successfully"
// @Failure      400   {string} string  "Bad request"
// @Failure      500   {string} string  "Internal server error"
// @Router       /songs [post]
func AddSongHandler(w http.ResponseWriter, r *http.Request) {
	SongsHandler(w, r)
}

// @Summary      Update song data
// @Description  Update song information in the database.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        song       body    SongData    true  "Song data to update"
// @Success      200       {string} string  "Song updated"
// @Failure      400       {string} string  "Bad request"
// @Failure      404       {string} string  "Song not found"
// @Failure      500       {string} string  "Internal server error"
// @Router       /songs [patch]
func UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	SongsHandler(w, r)
}

// @Summary      Delete a song
// @Description  Delete a song from the database.
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        song   query    string  true   "Song name"
// @Param        group  query    string  true   "Group name"
// @Success      200    {string} string  "Song deleted"
// @Failure      400    {string} string  "Bad request"
// @Failure      404    {string} string  "Song not found"
// @Failure      500    {string} string  "Internal server error"
// @Router       /songs [delete]
func DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	SongsHandler(w, r)
}

// @Summary      Get song lyrics
// @Description  Get the lyrics of a song by its name and group.
// @Tags         text
// @Accept       json
// @Produce      json
// @Param        song   query    string  true   "Song name"
// @Param        group  query    string  true   "Group name"
// @Param        verse  query    int     false  "Verse number"
// @Success      200    {string} string  "Song lyrics"
// @Failure      400    {string} string  "Bad request"
// @Failure      404    {string} string  "Song not found"
// @Failure      500    {string} string  "Internal server error"
// @Router       /text [get]
func GetSongTextHandler(w http.ResponseWriter, r *http.Request) {
	SongsHandler(w, r)
}
