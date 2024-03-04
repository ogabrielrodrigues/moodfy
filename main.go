package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/ogabrielrodrugues/moodfy/internal/artist"
	"github.com/ogabrielrodrugues/moodfy/internal/music"
	"github.com/ogabrielrodrugues/moodfy/internal/style"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	origin := os.Getenv("ORIGIN")
	if origin == "" {
		origin = "*"
	}

	db_url := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mh := music.Handler(db, origin)
	ah := artist.Handler(db, origin)
	sh := style.Handler(db, origin)

	mux.HandleFunc("POST /music", mh.CreateMusic)
	mux.HandleFunc("GET /music", mh.ListMusic)
	mux.HandleFunc("POST /artist", ah.CreateArtist)
	mux.HandleFunc("GET /artist", ah.ListArtist)
	mux.HandleFunc("POST /style", sh.CreateStyle)
	mux.HandleFunc("GET /style", sh.ListStyle)

	http.ListenAndServe(":8080", mux)
}
