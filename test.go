package main

import (
	"fmt"
	"music/internal/database"
	"net/url"
)

func main() {
	songs, _ := database.ListSongs(url.Values{"group": []string{"Portishead"}})
	fmt.Println(songs)
}
