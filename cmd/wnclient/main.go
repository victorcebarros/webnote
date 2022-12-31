package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"webnote"
)

func addNote(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s addr note", args[0])
	}

	addr := args[1]
	notebody := strings.Join(args[2:], " ")
	author, _ := os.Hostname()

	note := webnote.Note{
		Author: webnote.ClientID(author),
		Note:   notebody,
	}

	buf, err := json.Marshal(note)

	if err != nil {
		return err
	}

	resp, err := http.Post(addr+"/add",
		"application/json",
		bytes.NewReader(buf),
	)

	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}

func fetchAllNotes(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("usage: %s addr cid", args[0])
	}

	addr := args[1]
	cid := args[2]

	resp, err := http.Get(addr + "/fetchall?cid=" + cid)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	notes := []webnote.Note{}

	err = json.Unmarshal(body, &notes)

	if err != nil {
		return err
	}

	for _, note := range notes {
		fmt.Println(note)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s cmd ...\n", os.Args[0])
		os.Exit(1)
	}

	err := error(nil)

	switch os.Args[1] {
	case "add":
		err = addNote(os.Args[1:])
	case "fetch":
		err = fetchAllNotes(os.Args[1:])
	case "help":
		fmt.Println("subcommands available: add fetch help")
	default:
		err = fmt.Errorf("invalid subcommand: %s", os.Args[1])
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: %v\n", os.Args[0], err)
	}
}
