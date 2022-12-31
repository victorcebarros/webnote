// this package is **not** concurrency-safe

package webnote

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type ClientID string
type NoteID string

type Note struct {
	Note   string
	Author ClientID
}

func (n Note) String() string {
	return fmt.Sprintf("note by (%s):\n%s", n.Author, n.Note)
}

type NoteMap map[NoteID]Note
type SeenNotes map[NoteID]struct{}
type SeenMap map[ClientID]SeenNotes

type NoteDB struct {
	Notes NoteMap
	Seen  SeenMap
}

func NewNoteDB() *NoteDB {
	return &NoteDB{Notes: make(NoteMap), Seen: make(SeenMap)}
}

func (n *NoteDB) add(note Note) {
	id := uuid.New().String()

	for _, exists := n.Notes[NoteID(id)]; exists; {
		id = uuid.New().String()
	}

	n.Notes[NoteID(id)] = note
}

func (n *NoteDB) fetchAll(cid ClientID) []Note {
	clientSeen, ok := n.Seen[cid]

	if !ok {
		n.Seen[cid] = make(SeenNotes)
		clientSeen = n.Seen[cid]
	}

	notes := []Note{}

	for nid, note := range n.Notes {
		if _, exists := clientSeen[nid]; exists {
			continue
		}

		notes = append(notes, note)
		clientSeen[nid] = struct{}{}
	}

	return notes
}

func (n *NoteDB) Add(w http.ResponseWriter, req *http.Request) {
	buf, err := io.ReadAll(req.Body)

	if err != nil {
		log.Printf("%s: error: %s\n", os.Args[0], err)
		return
	}

	defer req.Body.Close()

	note := Note{}

	err = json.Unmarshal(buf, &note)

	if note.Author == "" {
		note.Author = ClientID(uuid.New().String())
	}

	if err != nil {
		log.Printf("%s: error: %s\n", os.Args[0], err)
		return
	}

	n.add(note)
}

func (n *NoteDB) FetchAll(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	values := req.URL.Query()

	cid := ClientID(values.Get("cid"))

	if cid == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found\n"))
		return
	}

	notes := n.fetchAll(cid)

	buf, err := json.Marshal(notes)

	if err != nil {
		log.Printf("%s: error: %s\n", os.Args[0], err)
	}

	w.Write(buf)
}
