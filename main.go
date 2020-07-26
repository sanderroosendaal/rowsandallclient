package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

var (
	config = oauth2.Config{
		ClientID:     "QrpI3SQW4tijWPyfboVuvqnmgE3WAMnnvAuvTN5r",
		ClientSecret: "RJoMU32qfWIqDLR9DU581h0mZhOP9sJTqcRE5flxaP5GzfZbovVTr8eYquFH5o2T8F7mwYXzkJhLTUeOXU2RmCWDsxkI2N1y4UsSBdhlpQ0sIRmN0NmaWUuBWNuENHSA",
		// Scopes:       []string{"read,write"},
		RedirectURL: "http://localhost:9094/oauth2",
		// This points to our Authorization Server
		// if our Client ID and Client Secret are valid
		// it will attempt to authorize our user
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:8000/rowers/o/authorize/",
			TokenURL: "http://localhost:8000/rowers/o/token/",
		},
	}
)

var stoken string

// WorkoutBody defines json for workout
type WorkoutBody struct {
	Name        string `json:"name"`
	Date        string `json:"date"`
	WorkoutType string `json:"workouttype"`
	StartTime   string `json:"starttime"`
	Distance    int64  `json:"distance"`
	Duration    string `json:"duration"`
}

// http://localhost:8000/rowers/o/authorize?client_id=QrpI3SQW4tijWPyfboVuvqnmgE3WAMnnvAuvTN5r&state=random_state_string&response_type=code

// HomePage serves home page
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Homepage Hit!")
	u := config.AuthCodeURL("xyz")
	http.Redirect(w, r, u, http.StatusFound)
}

// Workouts get workouts
func Workouts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Requesting Workouts")
	url := "http://localhost:8000/rowers/api/workouts/"
	// Create a Bearer string by appending string access token
	var bearer = fmt.Sprintf("Bearer %s", stoken)
	log.Println(bearer)

	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)

	// add authorization header to the req
	req.Header.Set("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", body)
}

// AddWorkout post a workout
func AddWorkout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Posting a workout")

	url := "http://localhost:8000/rowers/api/workouts/"
	var bearer = fmt.Sprintf("Bearer %s", stoken)
	log.Println(bearer)

	bodystruct := WorkoutBody{
		Name:        "From Go Client",
		Date:        "2020-07-25",
		WorkoutType: "water",
		StartTime:   "09:05:00",
		Distance:    12124,
		Duration:    "01:05:23",
	}

	bodyjson, err := json.Marshal(bodystruct)
	fmt.Printf("%s\n", []byte(bodyjson))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyjson)))
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response. \n[ERRO] -", err)
	}
	fmt.Printf("Respons code %v\n", resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", body)

}

// Authorize to do authorization
func Authorize(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	state := r.Form.Get("state")
	if state != "xyz" {
		http.Error(w, "State invalid", http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := config.Exchange(context.Background(), code)
	stoken = token.AccessToken
	log.Println(stoken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(*token)
}

func main() {
	stoken = "aap"
	// 1 - We attempt to hit our Homepage route
	// if we attempt to hit this unauthenticated, it
	// will automatically redirect to our Auth
	// server and prompt for login credentials
	http.HandleFunc("/", HomePage)

	// 2 - This displays our state, code and
	// token and expiry time that we get back
	// from our Authorization server
	http.HandleFunc("/oauth2", Authorize)

	http.HandleFunc("/workouts", Workouts)

	http.HandleFunc("/workout", AddWorkout)

	// 3 - We start up our Client on port 9094
	log.Println("Client is running at 9094 port.")
	log.Fatal(http.ListenAndServe(":9094", nil))
}
