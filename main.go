package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	Id       string `json:"Id"`
  Username string `json:"Username"`
  Email    string `json:"Email"`
  Password string `json:"Password"`
}
type Bottle struct {
	Id      string  `json:"Id"`
	Title   string  `json:"Title"`
	Content string  `json:"Content"`
}
type Bottle_Tag struct {
  Id       string `json:"Id"`
  BottleId string `json:"BottleId"`
  Tag      string `json:"tag"`
}
type Incident struct {
	Id       string `json:"Id"`
  ClientId string `json:"ClientId"`
	Content  string `json:"Content"`
	Location string `json:"Location"`
	Time     string `json:"Time"`
}
type Incident_Type struct {
  Id         string `json:"Id"`
  IncidentId string `json:"IncidentId"`
  Type       string `json:"Type"`
}
type Message struct {
  Id         string `json:"Id"`
  ClientId   string `json:"ClientId"`
  Time       string `json:"Time"` 
  Content    string `json:"Content"`
  ToClientId string `json:"ToClientId"`
}

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var bottlesIdCounter int
var incidentsIdCounter int
var messagesIdCounter int
var db *gorm.DB

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllBottles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	w.Header().Set("Content-Type", "application/json")

  var bottles []Bottle
	db.Exec("USE ocean;")
  db.Raw("SELECT * FROM bottle;").Scan(&bottles);
	json.NewEncoder(w).Encode(bottles)
}

func returnBottleById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnBottleById")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

  var bottle Bottle
	db.Exec("USE ocean;")
  db.Raw("SELECT * FROM bottle WHERE id=?;", params["id"]).Scan(&bottle);
	json.NewEncoder(w).Encode(bottle)
}

func returnRandomBottle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnRandomBottle")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	matchingBottles := []Bottle{}
	val, ok := params["tag"]
  db.Exec("USE ocean;")
	if ok {
    var suitableBottles []Bottle
    db.Raw("SELECT * FROM bottle WHERE Id IN (SELECT bottleid FROM bottle_tag WHERE tag=?);", val).Scan(&suitableBottles)
    matchingBottles = append(matchingBottles, suitableBottles...)
	} else {
    var bottles []Bottle
    db.Raw("SELECT * FROM bottle;").Scan(&bottles)
		json.NewEncoder(w).Encode(bottles[rand.Intn(len(bottles)+1)])
		return
	}

	min := 0
	max := len(matchingBottles)
	if max == min {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("404 - bro there's nothing there"))
		return
	}

	json.NewEncoder(w).Encode(matchingBottles[rand.Intn(max-min+1)])

}

func stringExists(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func deleteBottleById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteBottleById")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

  db.Exec("USE ocean;")
  db.Exec("DELETE FROM bottle WHERE id=?;", params["id"])
}

func createBottle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createBottle")
	w.Header().Set("Content-Type", "application/json")
	var bottle Bottle
	json.NewDecoder(r.Body).Decode(&bottle)
	bottle.Id = strconv.Itoa(bottlesIdCounter)
	bottlesIdCounter++
  db.Exec("USE ocean;")
  db.Exec("INSERT INTO bottle (Id,Title,Content) VALUES (?,?,?);", bottle.Id, bottle.Title, bottle.Content)
	json.NewEncoder(w).Encode(bottle)
}

func loginAuth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

  good := `{"Id": 1}`
  json.NewEncoder(w).Encode(good);
}

func createIncidentReport(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createIncidentReport")
	w.Header().Set("Content-Type", "application/json")
	var incident Incident
	json.NewDecoder(r.Body).Decode(&incident)
	incident.Id = strconv.Itoa(incidentsIdCounter)
	incidentsIdCounter++

  db.Exec("USE ocean;")
  db.Exec("INSERT INTO incident (Id,ClientId,Location,Content,Time) VALUES (?,?,?,?,?);", incident.Id, incident.ClientId, incident.Location, incident.Content, incident.Time)

	json.NewEncoder(w).Encode(incident)
}

func returnAllIncidents(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Endpoint Hit: returnAllIncidents")
  w.Header().Set("Content-Type", "application/json")
  var incidents []Incident
  db.Exec("USE ocean;")
  db.Raw("SELECT * FROM incident;").Scan(&incidents)

  json.NewEncoder(w).Encode(incidents)
}

func returnIncidentById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnIncidentById")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

  var incident Incident
	db.Exec("USE ocean;")
  db.Raw("SELECT * FROM incident WHERE id=?;", params["id"]).Scan(&incident);
	json.NewEncoder(w).Encode(incident)
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createMessage")
	w.Header().Set("Content-Type", "application/json")
	var message Message
	json.NewDecoder(r.Body).Decode(&message)
	message.Id = strconv.Itoa(messagesIdCounter)
	messagesIdCounter++
  db.Exec("USE ocean;")
  db.Exec("INSERT INTO message (Id,ClientId,Time,Content,ToClientId) VALUES (?,?,?,?,?);", message.Id, message.ClientId, message.Time, message.Content, message.ToClientId)
	json.NewEncoder(w).Encode(message)
}

func returnMessagesByClientId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnMessageByClientId")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

  var messages []Message
	db.Exec("USE ocean;")
  db.Raw("SELECT * FROM message WHERE clientid=?;", params["clientid"]).Scan(&messages);
  json.NewEncoder(w).Encode(messages)
}

func returnMessagesByToClientId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnMessageByToClientId")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

  var messages []Message
	db.Exec("USE ocean;")
  db.Raw("SELECT * FROM message WHERE toclientid=?;", params["toclientid"]).Scan(&messages);
  json.NewEncoder(w).Encode(messages)
}

func returnMessageById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnMessageById")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

  var message Message
	db.Exec("USE ocean;")
  db.Raw("SELECT * FROM message WHERE id=?;", params["id"]).Scan(&message)
  json.NewEncoder(w).Encode(message)
}

func getTags(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Endpoint Hit: getTags")
  w.Header().set("Content-Type", "application/json")
  params := mux.Vars(r)

  var tags []Bottle_Tag
  db.Exec("USE ocean;")
  db.Raw("SELECT * FROM bottle_tag WHERE clientid=?;", params["clientid"]).Scan(&tags)
  json.NewEncoder(w).Encode(tags)
}
func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage).Methods("GET")
	myRouter.HandleFunc("/bottles/all", returnAllBottles).Methods("GET")
	myRouter.HandleFunc("/bottles/{id}", returnBottleById).Methods("GET")
	myRouter.HandleFunc("/bottles/getRandom/", returnRandomBottle).Methods("GET")
	myRouter.HandleFunc("/bottles/getRandom/{tag}", returnRandomBottle).Methods("GET")
  myRouter.HandleFunc("/bottles/getTag/{clientid}", returnTags).Methods("GET")
	myRouter.HandleFunc("/login", loginAuth).Methods("POST")
	myRouter.HandleFunc("/bottles/{id}", deleteBottleById).Methods("DELETE")
	myRouter.HandleFunc("/bottles", createBottle).Methods("POST")
	myRouter.HandleFunc("/incidents", createIncidentReport).Methods("POST")
  myRouter.HandleFunc("/incidents/all", returnAllIncidents).Methods("GET") 
  myRouter.HandleFunc("/incidents/{id}", returnIncidentById).Methods("GET") 
  myRouter.HandleFunc("/messages", createMessage).Methods("POST")
  myRouter.HandleFunc("/messages/{id}", returnMessageById).Methods("GET")
  myRouter.HandleFunc("/messages/from/{clientid}", returnMessagesByClientId).Methods("GET")
  myRouter.HandleFunc("/messages/to/{toclientid}", returnMessagesByToClientId).Methods("GET")

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	fmt.Println("starting to connect to cockroachdb cluster...")
	// Connect to the "ocean" database as the "kaiser" user.
	// Read in connection string
	scanner := bufio.NewScanner(os.Stdin)
	log.Println("Enter a connection string: ")
	scanner.Scan()
	connstring := os.ExpandEnv(scanner.Text())

	// Connect to the "ocean" database
	var err error
	db, err = gorm.Open(postgres.Open(connstring), &gorm.Config{})
	if err != nil {
		log.Fatal("error configuring the database: ", err)
	}

	log.Println("Hey! You successfully connected to your CockroachDB cluster.")
	fmt.Println("started localhost @ 127.0.0.1:10000")
	bottlesIdCounter = 1
  incidentsIdCounter = 1
  messagesIdCounter = 1
	handleRequests()
}
