package main

import (
	"log"
	"net/http"
	"webnote"
)

func main() {
	db := webnote.NewNoteDB()

	http.HandleFunc("/fetchall", db.FetchAll)
	http.HandleFunc("/add", db.Add)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
