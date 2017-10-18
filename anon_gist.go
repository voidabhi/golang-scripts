
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//flags
var (
	publicFlag      bool
	descriptionFlag string
	responseObj     map[string]interface{}
)

//types
type GistFile struct {
	Content string `json:"content"`
}

type Gist struct {
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]GistFile `json:"files"`
}

func main() {
	flag.BoolVar(&publicFlag, "p", true, "Set to false for private gist.")
	flag.StringVar(&descriptionFlag, "d", "", "Description for gist.")
	flag.Parse()

	filenames := flag.Args()
	if len(filenames) == 0 {
		log.Fatal("Error: No files specified.")
	}

	fmt.Println(filenames)
	files := map[string]GistFile{}


	for _, filename := range filenames {
		fmt.Println("Checking file:", filename)
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal("File Error: ", err)
		}
		files[filename] = GistFile{string(content)}
	}

	if descriptionFlag == "" {
		descriptionFlag = strings.Join(filenames, ", ")
	}

	gist := Gist{
		descriptionFlag,
		publicFlag,
		files,
	}

	b, err := json.Marshal(gist)
	if err != nil {
		log.Fatal("JSON Error: ", err)
	}
	fmt.Println("OK")

	br := bytes.NewBuffer(b)
	fmt.Println("Uploading.")
	resp, err := http.Post("https://api.github.com/gists", "application/json", br)
	if err != nil {
		log.Fatal("HTTP Error: ", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&responseObj)
	if err != nil {
		log.Fatal("Response JSON Error: ", err)
	}

	fmt.Println("--- URL for gist ---")
	fmt.Println(responseObj["html_url"])
}
