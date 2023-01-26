package main

import (
	"bytes"
	"flag"
	"context"
	"encoding/json"
	"os"
	"io"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/sanderroosendaal/gorow"

	"golang.org/x/oauth2"
	"moul.io/http2curl"
)

var (
	config = oauth2.Config{
		ClientID: "YCQ1hVjdRECtin2NN75HHPjlni63uAHUWbgNEcGT", // local
		ClientSecret: "LKiXu7dL1biHtM84kMtYZnGbcbokWJVyA0tVPc9laRWJjpzFaNwTZmYruN7iGH6wVYC1kiH9HbxekWGL59XLBY3AVew5R3xNKgPfwq8G1LNglrriyBRPZXKFQTdjxOCx", // local
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
	config_dev = oauth2.Config{
		ClientID: "d3EXkobGFIdWOWkPwgXyth8So3CdVTEQLZGVw8qh", // dev
		ClientSecret: "wCg2brL8wDXhXcUTRDqoWm2h4naw8WORNUDHHO5dQ1AqBR62w8qsCmYrucMJHGWJ36OXJZqhZwyAy4to92ACdlNo6jayipwG98eP5hdsc213zPSleJgOFilYAECvtjE2", // dev
		// Scopes:       []string{"read,write"},
		RedirectURL: "http://localhost:9094/oauth2",
		// This points to our Authorization Server
		// if our Client ID and Client Secret are valid
		// it will attempt to authorize our user
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://dev.rowsandall.com/rowers/o/authorize/",
			TokenURL: "https://dev.rowsandall.com/rowers/o/token/",
		},
	}
	config_prod = oauth2.Config{
		ClientID: "SpecU3qiWht3mVk5VAXYPsaryW1PnsyuAecCTWcZ", // prod
		ClientSecret: "Oh9lUhw5YeTekFy5fYT4Qj2Je4qqBcduMzvOQgA2JWlJw3Pom575KmFGPR6GOdf43vxRmOgdxSPuzqS9N2XMvNVyacgQjtHULsf96t3ouqkacIBlZGPT8jh1pA5ZSV1M", // prod
		// Scopes:       []string{"read,write"},
		RedirectURL: "http://localhost:9094/oauth2",
		// This points to our Authorization Server
		// if our Client ID and Client Secret are valid
		// it will attempt to authorize our user
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://rowsandall.com/rowers/o/authorize/",
			TokenURL: "https://rowsandall.com/rowers/o/token/",
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

// StrokeRecord defines json for strokes data
type StrokeRecord struct {
	Distance  float64 `json:"distance"`
	Power     float64 `json:"power"`
	Heartrate int32   `json:"hr"`
	Pace      float64 `json:"pace"`
	Time      float64 `json:"time"`
	Spm       float64 `json:"spm"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Strokes defines json for strokes data
type Strokes struct {
	Data []StrokeRecord `json:"data"`
}

// Workout for API v3
type Workout struct {
	Name        string `json:"name"`
	WorkoutType string `json:"workouttype"`
	BoatType string `json:"boattype"`
	Notes string `json:"notes"`
	StartDateTime   string `json:"startdatetime"`
	Distance    int64  `json:"distance"`
	ElapsedTime    int64 `json:"elapsedTime"`
	Duration    string `json:"duration"`
	Strokes     Strokes `json:"strokes"`
}

// http://localhost:8000/rowers/o/authorize?client_id=QrpI3SQW4tijWPyfboVuvqnmgE3WAMnnvAuvTN5r&state=random_state_string&response_type=code

// HomePage serves home page
func HomePage(w http.ResponseWriter, r *http.Request) {
	if verbose{
		log.Println("Homepage Hit!")
	}
	u := config.AuthCodeURL("xyz")
	if instance == "dev" {
		u = config_dev.AuthCodeURL("xyz")
	}
	if instance == "prod" {
		u = config_prod.AuthCodeURL("xyz")
	} 
	http.Redirect(w, r, u, http.StatusFound)
}

// Workouts get workouts
func Workouts(w http.ResponseWriter, r *http.Request) {
	if verbose {
		log.Println("Requesting Workouts")
	}
	url := "http://localhost:8000/rowers/api/workouts/"
	if instance == "dev" {
		url = "https://dev.rowsandall.com/rowers/api/workouts/"		
	}
	if instance == "prod" {
		url = "https://rowsandall.com/rowers/api/workouts/"
	}

	// Create a Bearer string by appending string access token
	var bearer = fmt.Sprintf("Bearer %s", stoken)
	if verbose {
		log.Println(bearer)
	}
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
	if verbose{
		fmt.Printf("Response code %v\n", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", body)
}

// StrokeData v3
func StrokeDataV3(w http.ResponseWriter, r *http.Request) {
	url := "http://localhost:8000/rowers/api/v3/workouts/"
	if instance == "dev" {
		url = "https://dev.rowsandall.com/rowers/api/v3/workouts/"		
	}
	if instance == "prod" {
		url = "https://rowsandall.com/rowers/api/v3/workouts/"
	}

	var bearer = fmt.Sprintf("Bearer %s", stoken)

	file, _ := ioutil.ReadFile("teststrokes2.json")
	//file, _ := ioutil.ReadFile("RC-Upload-RAA.json")
	data := Workout{}
	_ = json.Unmarshal([]byte(file), &data)

	bodyjson, _ := json.Marshal(data)
	if verbose {
		fmt.Printf("%s\n", []byte(bodyjson))
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyjson)))
	// req, _ := http.NewRequest("POST", url, bytes.NewBuffer(file))
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	command, _ := http2curl.GetCurlCommand(req)
	if verbose {
		fmt.Println(command)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response. \n[ERRO] -", err)
	}
	if verbose {
		fmt.Printf("Response code %v\n", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", body)
}

// StrokeData Add strokedata
func StrokeData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form.Get("id")
	if verbose {
		fmt.Println(id)
	}
	url := fmt.Sprintf("http://localhost:8000/rowers/api/v2/workouts/%s/strokedata/", id)
	if instance == "dev" {
		url = fmt.Sprintf("https://dev.rowsandall.com/rowers/api/v2/workouts/%s/strokedata/", id)
	}
	if instance == "prod" {
		url = fmt.Sprintf("https://rowsandall.com/rowers/api/v2/workouts/%s/strokedata/", id)
	}

	var bearer = fmt.Sprintf("Bearer %s", stoken)

	file, _ := ioutil.ReadFile("teststrokes.json")
	//file, _ := ioutil.ReadFile("RC-Upload-RAA.json")
	data := Strokes{}
	_ = json.Unmarshal([]byte(file), &data)

	bodyjson, _ := json.Marshal(data)
	if verbose {
		fmt.Printf("%s\n", []byte(bodyjson))
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyjson)))
	// req, _ := http.NewRequest("POST", url, bytes.NewBuffer(file))
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	command, _ := http2curl.GetCurlCommand(req)
	if verbose {
		fmt.Println(command)		
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response. \n[ERRO] -", err)
	}
	if verbose {
		fmt.Printf("Response code %v\n", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", body)
}

// AddWorkout post a workout
func AddWorkout(w http.ResponseWriter, r *http.Request) {
	if verbose {
		fmt.Println("Posting a workout")
	}

	url := "http://localhost:8000/rowers/api/workouts/"
	if instance == "dev" {
		url = "https://dev.rowsandall.com/rowers/api/workouts/"		
	}
	if instance == "prod" {
		url = "https://rowsandall.com/rowers/api/workouts/"
	}

	var bearer = fmt.Sprintf("Bearer %s", stoken)
	if verbose {
		log.Println(bearer)		
	}


	bodystruct := WorkoutBody{
		Name:        "From Go Client",
		Date:        "2020-07-25",
		WorkoutType: "water",
		StartTime:   "09:05:00",
		Distance:    12124,
		Duration:    "01:05:23",
	}

	bodyjson, _ := json.Marshal(bodystruct)
	if verbose {
		fmt.Printf("%s\n", []byte(bodyjson))
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyjson)))
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	command, _ := http2curl.GetCurlCommand(req)
	if verbose {
		fmt.Println(command)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response. \n[ERRO] -", err)
	}
	if verbose {
		fmt.Printf("Response code %v\n", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", body)

}

func WorkoutForm(w http.ResponseWriter, r *http.Request) {
	url := "http://localhost:8000/rowers/api/v3/workouts/"
	if instance == "dev" {
		url = "https://dev.rowsandall.com/rowers/api/v3/workouts/"		
	}
	if instance == "prod" {
		url = "https://rowsandall.com/rowers/api/v3/workouts/"		
	}


	var bearer = fmt.Sprintf("Bearer %s", stoken)
	switch r.Method {
	case "GET": {
		http.ServeFile(w, r, "static/form.html")
	}
	case "POST": {
		if err := r.ParseMultipartForm(0); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "Post from website! r.PostForm = %v\n", r.PostForm)
		name := r.FormValue("name")
		workouttype := string(r.FormValue("workouttype"))
		boattype := string(r.FormValue("boattype"))
		notes := string(r.FormValue("notes"))
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			fmt.Fprintf(w, "Error Retrieving the File err: %v", err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "Uploaded File: %+v\n", handler.Filename)
		fmt.Fprintf(w, "File Size: %+v\n", handler.Size)
		fmt.Fprintf(w, "MIME Header: %+v\n", handler.Header)
		fmt.Fprintf(w, "Name = %s\n", name)
		fmt.Fprintf(w, "Workout Type = %s\n", workouttype)
		fmt.Fprintf(w, "Boat Type = %s\n", boattype)
		fmt.Fprintf(w, "Notes = %s\n", notes)

		// Create file
		dst, err := os.Create(handler.Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Copy the uploaded file to the created file on the filesystem
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Successfully Uploaded File\n")
		
		strokes, err := gorow.ReadCSV(handler.Filename)
		if err != nil {
			fmt.Fprintf(w, "Error Reading the CSV File err: %v", err)
			return
		}
		// we're good
		workout := Workout{}
		workout.Name = name
		workout.WorkoutType = workouttype
		workout.BoatType = boattype
		workout.Notes = notes
		workout.StartDateTime = time.Now().Format(time.RFC3339)
		workout.Distance = int64(strokes[len(strokes)-1].Distance)
		workout.ElapsedTime = int64((strokes[len(strokes)-1].Timestamp-strokes[0].Timestamp)*1000.)
		record := StrokeRecord{}
		data := Strokes{}
		for i := range strokes {
			record.Distance = strokes[i].Distance
			record.Power = strokes[i].Power
			record.Heartrate = int32(strokes[i].Hr)
			record.Pace = 1000.*strokes[i].Pace
			record.Time = (strokes[i].Timestamp-strokes[0].Timestamp)*1000.
			record.Spm = strokes[i].Spm
			record.Latitude = strokes[i].Latitude
			record.Longitude = strokes[i].Longitude
			data.Data = append(data.Data,record)
		}
		workout.Strokes = data

		bodyjson, _ := json.Marshal(workout)
		if verbose {
			fmt.Printf("%s\n", []byte(bodyjson))
		}

		req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(bodyjson)))
		// req, _ := http.NewRequest("POST", url, bytes.NewBuffer(file))
		req.Header.Set("Authorization", bearer)
		req.Header.Add("Content-Type", "application/json")

		command, _ := http2curl.GetCurlCommand(req)
		if verbose {
			fmt.Println(command)
		}
		
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response. \n[ERRO] -", err)
		}
		if verbose {
			fmt.Printf("Response code %v\n", resp.StatusCode)
		}
		
		body, _ := ioutil.ReadAll(resp.Body)

		fmt.Fprintf(w, "%s", body)

		// delete file
		err = os.Remove(handler.Filename)
		if err != nil{
			fmt.Fprintf(w, "Error removing temp file: %v", err)
			
		}

	}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
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

	u_config := config
	if instance == "dev" {
		u_config = config_dev
	}
	if instance == "prod" {
		u_config = config_prod
	}
	token, err := u_config.Exchange(context.Background(), code)
	stoken = token.AccessToken
	refreshtoken = token.RefreshToken

	if verbose {
		log.Println(stoken)
		log.Println(refreshtoken)		
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(*token)
}

var instance = "local"
var verbose = false

func main() {
	stoken = "aap"
	refreshtoken = "noot"
	// flags
	flag.BoolVar(&(verbose), "v", false, "use -v to set verbose")
	flag.StringVar(&(instance), "i", "local", "use -i instance to set instance (local, dev, prod)")
	flag.Parse()
	if verbose {
		log.Println(instance)
		log.Println(time.Now().Format(time.RFC3339))		
	}

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

	http.HandleFunc("/form", WorkoutForm)

	// 3 - We start up our Client on port 9094
	log.Println("Client is running at 9094 port.")
	log.Println("Endpoints")
	log.Println("/")
	log.Println("/oauth2")
	log.Println("/workouts")
	log.Println("/workout")
	log.Println("/strokedata")
	log.Println("/strokedatav3")
	log.Println("/form")
	log.Fatal(http.ListenAndServe(":9094", nil))
}
