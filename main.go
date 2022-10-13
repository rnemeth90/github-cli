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
	token       = flag.String("token", "", "github auth token")
	owner       = flag.String("owner", "", "the repo owner")
	repoName    = flag.String("repoName", "", "github repo")
	title       = flag.String("title", "", "title for new issue")
	body        = flag.String("body", "", "body for new issue")
	createRepo  = flag.Bool("createRepo", false, "Create a repo")
	getrepos    = flag.Bool("getrepos", false, "Get all repos for a user")
	isPrivate   = flag.Bool("isPrivate", false, "Is the new repo private?")
	description = flag.String("description", "", "Description for the repo")
	help        = flag.Bool("help", false, "Print help")
)

// NewIssue - specify data fields for new github issue submission
type NewIssue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
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

	if *help {
		flag.PrintDefaults()
		return
	}

	if *createRepo {
		makeRepo(*repoName, *description, *token, *isPrivate)
	}

	if *getrepos {
		repos, err := getRepos(*owner, *token)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(PrettyPrint(repos))
	}
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

// createRepo creates a repo for an authenticated user
func makeRepo(name, description, token string, private bool) error {
	url := "https://api.github.com/user/repos"

	if name == "" || token == "" {
		fmt.Println("You must specify a repo name and token")
		os.Exit(2)
	}

	fmt.Printf("Creating repo: %s\n", name)

	repoData := Repository{
		Name:        name,
		Description: description,
		Private:     private,
	}

	jsonData, err := json.Marshal(repoData)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// we should return an error from this method, not this:
	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Response code is %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		//print body as it may contain hints in case of errors
		fmt.Println(string(body))
		log.Fatal(err)
	}
	return nil
}

// getRepos lists repos for an authenticated user
func getRepos(owner, token string) ([]Repository, error) {
	url := "https://api.github.com/user/repos"

	if owner == "" || token == "" {
		fmt.Println("You must specify an owner and token")
		os.Exit(2)
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var jsonData []Repository
	err = json.Unmarshal([]byte(body), &jsonData) // here!
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response code is %d\n", resp.StatusCode)
	}
	return jsonData, nil
}

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
