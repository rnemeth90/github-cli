package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const apiBaseUrl string = `https://api.github.com/repos/`

var (
	token = flag.String("token", "", "github auth token")
	owner = flag.String("owner", "", "the repo owner")
	repo  = flag.String("repo", "", "github repo")
	title = flag.String("title", "", "title for new issue")
	body  = flag.String("body", "", "body for new issue")
)

// NewIssue - specify data fields for new github issue submission
type NewIssue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func usage() {
	flag.PrintDefaults()
	return
}

func init() {
	flag.Usage = usage
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 0 {
		flag.Usage()
		os.Exit(2)
	}

	// createIssue(*repo, *owner, *title, *body, *token)
	createRepo(*owner, *repo)
}

func createIssue(repo, owner, title, body, token string) {
	apiURL := apiBaseUrl + owner + "/" + repo + "/issues"
	//title is the only required field
	issueData := NewIssue{Title: title, Body: body}
	//make it json
	jsonData, _ := json.Marshal(issueData)
	//creating client to set custom headers for Authorization
	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Response code is is %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		//print body as it may contain hints in case of errors
		fmt.Println(string(body))
		log.Fatal(err)
	}
}

//createRepo
func createRepo(owner, name string) {
	fmt.Printf("Creating repo: %s/%s\n", owner, name)
}
