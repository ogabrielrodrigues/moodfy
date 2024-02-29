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
	db_url := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mh := music.Handler(db)
	ah := artist.Handler(db)
	sh := style.Handler(db)

	mux.HandleFunc("POST /music", mh.CreateMusic)
	mux.HandleFunc("GET /music", mh.ListMusic)
	mux.HandleFunc("POST /artist", ah.CreateArtist)
	mux.HandleFunc("GET /artist", ah.ListArtist)
	mux.HandleFunc("POST /style", sh.CreateStyle)
	mux.HandleFunc("GET /style", sh.ListStyle)

	http.ListenAndServe(":8080", mux)
}
