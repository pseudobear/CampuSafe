package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Bottle struct {
	Id      string   `json:"Id"`
	Title   string   `json:"Title"`
	Tag     []string `json:"Tag"`
	Content string   `json:"Content"`
}

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var bottles []Bottle
var idCounter int

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllBottles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bottles)
}

func returnBottleById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnBottleById")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	json.NewEncoder(w).Encode(returnSpecificBottle(params["id"]))
}

func returnRandomBottle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnRandomBottle")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	matchingBottles := []Bottle{}
	val, ok := params["tag"]
	if ok {
		for _, item := range bottles {
			if stringExists(val, item.Tag) {
				matchingBottles = append(matchingBottles, item)
			}
		}
	} else {
		json.NewEncoder(w).Encode(bottles[rand.Intn(len(bottles)+1)])
		return
	}

	min := 0
	max := len(matchingBottles)
	if max == min {
		json.NewEncoder(w).Encode(returnSpecificBottle("0"))
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

func returnSpecificBottle(id string) Bottle {
	for _, item := range bottles {
		if item.Id == id {
			return item
		}
	}
	return Bottle{}
}

func deleteBottleById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteBottleById")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range bottles {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			bottles = append(bottles[:index], bottles[index+1:]...)
			return
		}
	}
}

func createBottle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createBottle")
	w.Header().Set("Content-Type", "application/json")
	var bottle Bottle
	json.NewDecoder(r.Body).Decode(&bottle)
	bottle.Id = strconv.Itoa(idCounter)
	idCounter++
	bottles = append(bottles, bottle)
	json.NewEncoder(w).Encode(bottle)
}

func updateBottle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: updateBottle")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range bottles {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			bottles = append(bottles[:index], bottles[index+1:]...)
			return
		}
	}
	var bottle Bottle
	json.NewDecoder(r.Body).Decode(&bottle)
	bottle.Id = params["id"]
	bottles = append(bottles, bottle)
	json.NewEncoder(w).Encode(bottle)
}

func loginAuth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage).Methods("GET")
	myRouter.HandleFunc("/bottles/all", returnAllBottles).Methods("GET")
	myRouter.HandleFunc("/bottles/{id}", returnBottleById).Methods("GET")
	myRouter.HandleFunc("bottles/getRandom/{tag}", returnRandomBottle).Methods("GET")
	myRouter.HandleFunc("/login", loginAuth).Methods("POST")
	myRouter.HandleFunc("/bottles/{id}", deleteBottleById).Methods("DELETE")
	myRouter.HandleFunc("/bottles", createBottle).Methods("POST")
	myRouter.HandleFunc("/bottles/{id}", updateBottle).Methods("PUT")
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	idCounter = 2
	bottles = []Bottle{
		Bottle{Id: "0", Title: "example", Tag: []string{"10"}, Content: "hello world"},
		Bottle{Id: "1", Title: "example", Tag: []string{"10", "20"}, Content: "hello world!"},
	}
	handleRequests()
}
