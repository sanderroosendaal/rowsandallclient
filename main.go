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
	"moul.io/http2curl"
)

var (
	config = oauth2.Config{
		ClientID: "YCQ1hVjdRECtin2NN75HHPjlni63uAHUWbgNEcGT", // local
		//ClientID: "d3EXkobGFIdWOWkPwgXyth8So3CdVTEQLZGVw8qh", // dev
		//ClientID: "SpecU3qiWht3mVk5VAXYPsaryW1PnsyuAecCTWcZ", // prod
		ClientSecret: "LKiXu7dL1biHtM84kMtYZnGbcbokWJVyA0tVPc9laRWJjpzFaNwTZmYruN7iGH6wVYC1kiH9HbxekWGL59XLBY3AVew5R3xNKgPfwq8G1LNglrriyBRPZXKFQTdjxOCx", // local
		//ClientSecret: "wCg2brL8wDXhXcUTRDqoWm2h4naw8WORNUDHHO5dQ1AqBR62w8qsCmYrucMJHGWJ36OXJZqhZwyAy4to92ACdlNo6jayipwG98eP5hdsc213zPSleJgOFilYAECvtjE2", // dev
		//ClientSecret: "Oh9lUhw5YeTekFy5fYT4Qj2Je4qqBcduMzvOQgA2JWlJw3Pom575KmFGPR6GOdf43vxRmOgdxSPuzqS9N2XMvNVyacgQjtHULsf96t3ouqkacIBlZGPT8jh1pA5ZSV1M", // prod
		// Scopes:       []string{"read,write"},
		RedirectURL: "http://localhost:9094/oauth2",
		// This points to our Authorization Server
		// if our Client ID and Client Secret are valid
		// it will attempt to authorize our user
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:8000/rowers/o/authorize/",
			TokenURL: "http://localhost:8000/rowers/o/token/",
			//AuthURL:  "https://dev.rowsandall.com/rowers/o/authorize/",
			//TokenURL: "https://dev.rowsandall.com/rowers/o/token/",
			//AuthURL:  "https://rowsandall.com/rowers/o/authorize/",
			//TokenURL: "https://rowsandall.com/rowers/o/token/",
		},
	}
)

var stoken string
var refreshtoken string

// WorkoutBody defines json for workout
type WorkoutBody struct {
	Name        string `json:"name"`
	Date        string `json:"date"`
	WorkoutType string `json:"workouttype"`
	StartTime   string `json:"starttime"`
	Distance    int64  `json:"distance"`
	Duration    string `json:"duration"`
}

// Strokes defines json for strokes data
type Strokes struct {
	Data []struct {
		Distance  float64 `json:"distance"`
		Power     float64 `json:"power"`
		Heartrate int32   `json:"hr"`
		Pace      float64 `json:"pace"`
		Time      float64 `json:"time"`
		Spm       float64 `json:"spm"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"data"`
}

// Workout for API v3
type Workout struct {
	Name        string `json:"title"`
	WorkoutType string `json:"workouttype"`
	BoatType string `json:"boattype"`
	Notes string `json:"notes"`
	StartDateTime   string `json:"startdatetime"`
	Distance    int64  `json:"totalDistance"`
	Duration    int64 `json:"elapsedTime"`
	Strokes     Strokes `json:"strokes"`
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
	//url := "https://dev.rowsandall.com/rowers/api/workouts/"
	//url := "https://rowsandall.com/rowers/api/workouts/"
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
	fmt.Printf("Response code %v\n", resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", body)
}

// StrokeData v3
func StrokeDataV3(w http.ResponseWriter, r *http.Request) {
	url := "http://localhost:8000/rowers/api/v3/workouts/"
	// url := "https://dev.rowsandall.com/rowers/api/v3/workouts/"
	// url := "https://rowsandall.com/rowers/api/v3/workouts/"

	var bearer = fmt.Sprintf("Bearer %s", stoken)

	file, _ := ioutil.ReadFile("teststrokes2.json")
	//file, _ := ioutil.ReadFile("RC-Upload-RAA.json")
	data := Workout{}
	_ = json.Unmarshal([]byte(file), &data)

	bodyjson, _ := json.Marshal(data)
	fmt.Printf("%s\n", []byte(bodyjson))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyjson)))
	// req, _ := http.NewRequest("POST", url, bytes.NewBuffer(file))
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	command, _ := http2curl.GetCurlCommand(req)
	fmt.Println(command)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response. \n[ERRO] -", err)
	}
	fmt.Printf("Response code %v\n", resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", body)
}

// StrokeData Add strokedata
func StrokeData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form.Get("id")
	fmt.Println(id)
	url := fmt.Sprintf("http://localhost:8000/rowers/api/v2/workouts/%s/strokedata/", id)
	// url := fmt.Sprintf("https://dev.rowsandall.com/rowers/api/v2/workouts/%s/strokedata/", id)
	//url := fmt.Sprintf("https://rowsandall.com/rowers/api/v2/workouts/%s/strokedata/", id)

	var bearer = fmt.Sprintf("Bearer %s", stoken)

	file, _ := ioutil.ReadFile("teststrokes.json")
	//file, _ := ioutil.ReadFile("RC-Upload-RAA.json")
	data := Strokes{}
	_ = json.Unmarshal([]byte(file), &data)

	bodyjson, _ := json.Marshal(data)
	fmt.Printf("%s\n", []byte(bodyjson))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyjson)))
	// req, _ := http.NewRequest("POST", url, bytes.NewBuffer(file))
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	command, _ := http2curl.GetCurlCommand(req)
	fmt.Println(command)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response. \n[ERRO] -", err)
	}
	fmt.Printf("Response code %v\n", resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", body)
}

// AddWorkout post a workout
func AddWorkout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Posting a workout")

	url := "http://localhost:8000/rowers/api/workouts/"
	//url := "https://dev.rowsandall.com/rowers/api/workouts/"
	//url := "https://rowsandall.com/rowers/api/workouts/"
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

	bodyjson, _ := json.Marshal(bodystruct)
	fmt.Printf("%s\n", []byte(bodyjson))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyjson)))
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	command, _ := http2curl.GetCurlCommand(req)
	fmt.Println(command)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response. \n[ERRO] -", err)
	}
	fmt.Printf("Response code %v\n", resp.StatusCode)

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
	refreshtoken = token.RefreshToken

	log.Println(stoken)
	log.Println(refreshtoken)
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
	refreshtoken = "noot"
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

	http.HandleFunc("/strokedata", StrokeData)

	http.HandleFunc("/strokedatav3", StrokeDataV3)

	// 3 - We start up our Client on port 9094
	log.Println("Client is running at 9094 port.")
	log.Println("Endpoints")
	log.Println("/")
	log.Println("/oauth2")
	log.Println("/workouts")
	log.Println("/workout")
	log.Println("/strokedata")
	log.Println("/strokedatav3")
	log.Fatal(http.ListenAndServe(":9094", nil))
}
