package main

import (
	"bytes"
	"flag"
	"context"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"os"
	"io"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/sanderroosendaal/gorow"
	"mime/multipart"
	"golang.org/x/oauth2"
	"moul.io/http2curl"
)

var (
	config = oauth2.Config{}
	apiworkouts_url string
	apicourses_url string
	apiv3_url string
	apiFIT_url string
	apiTCX_url string
	apirowingdata_url string
	apirowingdata_url_apikey string
	apistrokedata_url string
)

var Stoken string
var Refreshtoken string

type Config struct {
	ClientID string `yaml:"clientid"`
	ClientSecret string `yaml:"clientsecret"`
	RedirectURL string `yaml:redirecturl`
	ApiServer string `yaml:apiserver`
	//	ApiWorkouts string `yaml:apiworkouts`
	//	ApiV3 string `yaml:apiv3`
	// ApiStrokeData string `yamml:apistrokedata`
}

type User struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	APIKey string `yaml:"apikey"`
}
var user User

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



// HomePage serves home page
func HomePage(w http.ResponseWriter, r *http.Request) {
	if verbose{
		log.Println("Homepage Hit!")
	}
	u := config.AuthCodeURL("xyz")
	http.Redirect(w, r, u, http.StatusFound)
}

// func Courses gets courses
func Courses(w http.ResponseWriter, r *http.Request) {
	if verbose {
		log.Println("Requesting Courses")
	}
	url := apicourses_url
	// Create a Bearer string by appending string access token
	var bearer = fmt.Sprintf("Bearer %s", authkeys.Stoken)
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

// Workouts get workouts
func Workouts(w http.ResponseWriter, r *http.Request) {
	if verbose {
		log.Println("Requesting Workouts")
	}
	url := apiworkouts_url

	// Create a Bearer string by appending string access token
	var bearer = fmt.Sprintf("Bearer %s", authkeys.Stoken)
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

// Workouts get workouts
func WorkoutsAPI(w http.ResponseWriter, r *http.Request) {
	if verbose {
		log.Println("Requesting Workouts")
	}
	url := apiworkouts_url
	if verbose {
		log.Println(apiworkouts_url)
	}
	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating GET request")
		return
	}

	// Add Headers
	req.Header.Set("Authorization", user.APIKey)
	if verbose {
		fmt.Printf("APIKey %s\n", user.APIKey)
	}

	// get the CURL command
	command, _ := http2curl.GetCurlCommand(req)
	if verbose {
		fmt.Println(command)
	}
	

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

// StrokeData as FIT
func StrokeDataFIT(w http.ResponseWriter, r *http.Request) {
	url := apiFIT_url
	if instance == "dev" {
		url = "https://dev.rowsandall.com/rowers/api/FIT/workouts/"		
	}
	if instance == "prod" {
		url = "https://rowsandall.com/rowers/api/FIT/workouts/"
	}

	var bearer = fmt.Sprintf("Bearer %s", authkeys.Stoken)
	fitbody, _ := ioutil.ReadFile("fitdata.fit")

	req, _ := http.NewRequest("POST",url, bytes.NewBuffer(fitbody))
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Content-Type", "application/octet-stream")


	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response. \n[ERRO] -", err)
	}
	if verbose {
		fmt.Printf("Response code %v\n", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Fprintf(w, "%s", body)
}

// StrokeData as TCX
func StrokeDataTCX(w http.ResponseWriter, r *http.Request) {
	url := apiTCX_url
	if instance == "dev" {
		url = "https://dev.rowsandall.com/rowers/api/TCX/workouts/"		
	}
	if instance == "prod" {
		url = "https://rowsandall.com/rowers/api/TCX/workouts/"
	}

	var bearer = fmt.Sprintf("Bearer %s", authkeys.Stoken)

	xmlbody, _ := ioutil.ReadFile("tcxdata.tcx")

	if verbose {
		fmt.Printf("%s\n", []byte(xmlbody))
	}

	req, _ := http.NewRequest("POST",url, bytes.NewBuffer(xmlbody))
	req.Header.Set("Authorization", bearer)
	req.Header.Add("Content-Type", "application/xml")

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
	defer resp.Body.Close()
	fmt.Fprintf(w, "%s", body)
}

// StrokeData as RowingData data file
func StrokeDataRD(w http.ResponseWriter, r *http.Request) {
	url := apirowingdata_url
	if instance == "dev" {
		url = "https://dev.rowsandall.com/rowers/api/rowingdata/workouts/"		
	}
	if instance == "prod" {
		url = "https://rowsandall.com/rowers/api/rowingdata/workouts/"
	}

	// Open the CSV file
	csvfile, err := os.Open("testdata.csv")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer csvfile.Close()

	// Create a buffer to hold the multipart form data
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add the CSV file to the form
	part, err := writer.CreateFormFile("file", "data.csv")
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return
	}
	// Copy the file content into the form file part
	if _, err := io.Copy(part, csvfile); err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return
	}
	// Add additional form fields if needed
	if err := writer.WriteField("boattype", "1x"); err != nil {
		fmt.Printf("Error adding form field: %v\n", err)
		return
	}
	if err := writer.WriteField("workouttype", "rower"); err != nil {
		fmt.Printf("Error adding form field: %v\n", err)
		return
	}

	// Close the writer to finalize the form data
	if err := writer.Close(); err != nil {
		fmt.Printf("Error closing writer: %v\n", err)
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Add headers
	req.SetBasicAuth(user.Username, user.Password)
	req.Header.Set("Content-Type", writer.FormDataContentType())


	// get the CURL command
	// command, _ := http2curl.GetCurlCommand(req)
	//if verbose {
	//	fmt.Println(command)
	//}
	
	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}

	// Print response
	fmt.Printf("Response status: %s\n", resp.Status)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}
	fmt.Fprintf(w, "Response body: %s\n", string(responseBody))

}

// StrokeData as RowingData data file
func StrokeDataRDAPI(w http.ResponseWriter, r *http.Request) {
	url := apirowingdata_url_apikey
	if instance == "dev" {
		url = "https://dev.rowsandall.com/rowers/api/rowingdata/"		
	}
	if instance == "prod" {
		url = "https://rowsandall.com/rowers/api/rowingdata/"
	}

	// Open the CSV file
	csvfile, err := os.Open("testdata.csv")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer csvfile.Close()

	// Create a buffer to hold the multipart form data
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add the CSV file to the form
	part, err := writer.CreateFormFile("file", "data.csv")
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return
	}
	// Copy the file content into the form file part
	if _, err := io.Copy(part, csvfile); err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return
	}
	// Add additional form fields if needed
	if err := writer.WriteField("boattype", "1x"); err != nil {
		fmt.Printf("Error adding form field: %v\n", err)
		return
	}
	if err := writer.WriteField("workouttype", "rower"); err != nil {
		fmt.Printf("Error adding form field: %v\n", err)
		return
	}

	// Close the writer to finalize the form data
	if err := writer.Close(); err != nil {
		fmt.Printf("Error closing writer: %v\n", err)
		return
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, &body)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Add headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", user.APIKey)
	if verbose {
		fmt.Printf("APIKey %s\n", user.APIKey)
	}


	// get the CURL command
	command, _ := http2curl.GetCurlCommand(req)
	if verbose {
		fmt.Println(command)
	}

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}

	// Print response
	fmt.Printf("Response status: %s\n", resp.Status)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}
	fmt.Fprintf(w, "Response body: %s\n", string(responseBody))

}

// StrokeData v3
func StrokeDataV3(w http.ResponseWriter, r *http.Request) {
	url := apiv3_url
	if instance == "dev" {
		url = "https://dev.rowsandall.com/rowers/api/v3/workouts/"		
	}
	if instance == "prod" {
		url = "https://rowsandall.com/rowers/api/v3/workouts/"
	}

	var bearer = fmt.Sprintf("Bearer %s", authkeys.Stoken)

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
	url := fmt.Sprintf(apistrokedata_url, id)
	if instance == "dev" {
		url = fmt.Sprintf("https://dev.rowsandall.com/rowers/api/v2/workouts/%s/strokedata/", id)
	}
	if instance == "prod" {
		url = fmt.Sprintf("https://rowsandall.com/rowers/api/v2/workouts/%s/strokedata/", id)
	}

	var bearer = fmt.Sprintf("Bearer %s", authkeys.Stoken)

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

	url := apiworkouts_url
	if instance == "dev" {
		url = "https://dev.rowsandall.com/rowers/api/workouts/"		
	}
	if instance == "prod" {
		url = "https://rowsandall.com/rowers/api/workouts/"
	}

	var bearer = fmt.Sprintf("Bearer %s", authkeys.Stoken)
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
	url := apiv3_url
	if instance == "dev" {
		url = "https://dev.rowsandall.com/rowers/api/v3/workouts/"		
	}
	if instance == "prod" {
		url = "https://rowsandall.com/rowers/api/v3/workouts/"		
	}


	var bearer = fmt.Sprintf("Bearer %s", authkeys.Stoken)
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

// Refresh to do refresh
func Refresh(w http.ResponseWriter, r *http.Request) {
	log.Println("Refresh")
}

// Authorize to do authorization
func Authorize(w http.ResponseWriter, r *http.Request) {
        log.Println("Authorize")
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
	if err != nil {
		log.Println("Error on token exchange\n", err)
	}
	authkeys.Stoken = token.AccessToken
	authkeys.Refreshtoken = token.RefreshToken
	authkeys.Expiry = token.Expiry
	bodyyaml, err := yaml.Marshal(authkeys)
	fmt.Printf("%s \n", []byte(bodyyaml))

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	fileName := "tokens.yaml"
	err = ioutil.WriteFile(fileName, bodyyaml, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}

	if verbose {
		log.Println(Stoken)
		log.Println(Refreshtoken)		
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
var configfile = "config.yaml"
var authorized = false

type keys struct {
	Stoken string `yaml:"stoken"`
	Refreshtoken string `yaml:"clientsecret"`
	Expiry time.Time `yaml:"expiry"`
}

var authkeys = keys{}

func main() {
	// flags
	flag.BoolVar(&(verbose), "v", false, "use -v to set verbose")
	flag.StringVar(&(configfile), "c", "config.yaml", "use -c instance to set redirect config file")
	flag.BoolVar(&(authorized), "a", false, "use -a if you're already authorized")
	flag.Parse()
	if verbose {
		log.Println(instance)
		log.Println(time.Now().Format(time.RFC3339))		
	}
	if authorized {
		file, err := ioutil.ReadFile("tokens.yaml")
		if err == nil {
			err = yaml.Unmarshal([]byte(file), &authkeys)
			if err != nil {
				authorized = false
				log.Printf("tokens error %v\n", err)
			}
			if verbose {
				log.Printf("stoken %s\n",authkeys.Stoken)
				log.Printf("refreshtoken %s\n",authkeys.Refreshtoken)
			}
		} else {
			authorized = false
		}
	}

	if verbose {
		fmt.Println(configfile)
	}
	
	file, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		os.Exit(1)
	}
	
	newconfig := Config{}
	err = yaml.Unmarshal([]byte(file), &newconfig)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		os.Exit(1)
	}
	if verbose {
		bodyyaml, _ := yaml.Marshal(newconfig)
		fmt.Printf("%s \n", []byte(bodyyaml))
	}
	if newconfig.ClientSecret == "" {
		log.Printf("no ClientSecret in config file")
		os.Exit(1)
	}
	if newconfig.RedirectURL == "" {
		log.Printf("no RedirectURL in config file")
		os.Exit(1)
	}
	if newconfig.ApiServer == "" {
		log.Printf("no ApiServer in config file")
		os.Exit(1)
	}

	err = yaml.Unmarshal([]byte(file), &user)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		os.Exit(1)
	}
	if verbose {
		bodyyaml, _ := yaml.Marshal(user)
		fmt.Printf("%s \n", []byte(bodyyaml))
		fmt.Printf("username %s\n", user.Username)
		fmt.Printf("password %s\n", user.Password)
		fmt.Printf("APIKey %s\n", user.APIKey)
	}

	instance = "local"
	config = oauth2.Config{
		ClientID: newconfig.ClientID,
		ClientSecret: newconfig.ClientSecret,
		RedirectURL: newconfig.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL: newconfig.ApiServer+"/rowers/o/authorize/",
			TokenURL: newconfig.ApiServer+"/rowers/o/token/",
		},
	}
	apiworkouts_url = newconfig.ApiServer+"/rowers/api/workouts/"
	apicourses_url = newconfig.ApiServer+"/rowers/api/courses/kml/liked/"
	apiv3_url = newconfig.ApiServer+"/rowers/api/v3/workouts/"
	apiTCX_url = newconfig.ApiServer+"/rowers/api/TCX/workouts/"
	apiFIT_url = newconfig.ApiServer+"/rowers/api/FIT/workouts/"
	apirowingdata_url = newconfig.ApiServer+"/rowers/api/rowingdata/workouts/"
	apirowingdata_url_apikey = newconfig.ApiServer+"/rowers/api/rowingdata/"
	apistrokedata_url = newconfig.ApiServer+"/rowers/api/v2/workouts/%s/strokedata/"
	if verbose {
		log.Println(apiworkouts_url)
		log.Println(apiv3_url)
		log.Println(apistrokedata_url)
		log.Println(apiTCX_url)
		log.Println(apiFIT_url)
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

	http.HandleFunc("/courses", Courses)

	http.HandleFunc("/workouts", Workouts)
	http.HandleFunc("/workoutsAPI", WorkoutsAPI)

	http.HandleFunc("/workout", AddWorkout)

	http.HandleFunc("/strokedata", StrokeData)

	http.HandleFunc("/strokedatav3", StrokeDataV3)

	http.HandleFunc("/strokedataV3", StrokeDataV3)

	http.HandleFunc("/strokedataTCX", StrokeDataTCX)
	
	http.HandleFunc("/strokedataFIT", StrokeDataFIT)
	
	http.HandleFunc("/strokedataRD", StrokeDataRD)

	http.HandleFunc("/strokedataRDAPI", StrokeDataRDAPI)

	http.HandleFunc("/form", WorkoutForm)

	// 3 - We start up our Client on port 9094
	log.Println("Client is running at 9094 port.")
	log.Println("Endpoints")
	log.Println("/")
	log.Println("/oauth2")
	log.Println("/workouts")
	log.Println("/workoutsAPI")
	log.Println("/workout")
	log.Println("/strokedata")
	log.Println("/strokedatav3")
	log.Println("/strokedataTCX")
	log.Println("/strokedataFIT")
	log.Println("/strokedataRD")
	log.Println("/strokedataRDAPI")
	log.Println("/form")
	log.Fatal(http.ListenAndServe(":9094", nil))
}
