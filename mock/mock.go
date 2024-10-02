package mock

import (
	"encoding/json"
	"log"
	"net/http"
)

// Запускает сервер
func RunServer() {
	http.HandleFunc("/info", infoHandler)
	err := http.ListenAndServe("0.0.0.0:8181", nil)
	log.Fatal(err)
}

// Обработчик /info
func infoHandler(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	song := query.Get("song")
	group := query.Get("group")

	if song == "" || group == "" || len(query) != 2 {
		http.Error(writer, "Bad request", http.StatusBadRequest)
		return
	}

	if song == "Roads" && group == "Portishead" {
		songInfo := map[string]string{}
		songInfo["releaseDate"] = "22.08.1994"
		songInfo["link"] = "https://www.youtube.com/watch?v=Vg1jyL3cr60"
		songInfo["text"] = text

		writer.WriteHeader(200)
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(writer).Encode(songInfo)
		return
	} else {
		http.Error(writer, "Unknown song", http.StatusBadRequest)
		return
	}

}

var text string = `Oh
Can't anybody see
We've got a war to fight
Never find our way
Regardless of what they say

How can it feel this wrong?
From this moment
How can it feel this wrong?

Storm in the morning light
I feel, no more can I say
Frozen to myself
I got nobody on my side
And surely that ain't right
Surely that ain't right

Oh
Can't anybody see
We've got a war to fight
Never find our way
Regardless of what they say

How can it feel this wrong?
From this moment
How can it feel this wrong?

How can it feel this wrong?
From this moment
How can it feel this wrong?

Oh
Can't anybody see
We've got a war to fight
Never find our way
Regardless of what they say

How can it feel this wrong?
From this moment
How can it feel this wrong?`
