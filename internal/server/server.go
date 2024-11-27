package server

import (
	"encoding/json"
	"fmt"
	"io"
	"music/internal/services"
	"music/tools"
	"net/http"
	"strings"
)

// Запускает сервер
func RunServer() {
	serverAddr := tools.GetConfig().ServerAddr
	http.HandleFunc("/songs", songsHandler)
	http.HandleFunc("/text", textHandler)
	tools.Logger.Info(fmt.Sprintf("Starting server on %s", serverAddr))
	err := http.ListenAndServe(serverAddr, nil)
	tools.Logger.Fatal("Server is down: ", err)
}

// Обработчик /songs
func songsHandler(writer http.ResponseWriter, request *http.Request) {
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
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "Can't read request body", http.StatusBadRequest)
			return
		}
		defer request.Body.Close()

		var newSong map[string]string
		err = json.Unmarshal(body, &newSong)
		if err != nil {
			fmt.Println(err)
			http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		unexpectedParams, err := services.AddSong(newSong)
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
func textHandler(writer http.ResponseWriter, request *http.Request) {
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
