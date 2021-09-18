package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
)
type Bottle struct {
    Id string `json:"Id"`
    Title string `json:"Title"`
    Tag string `json:"Tag"`
    Content string `json:"Content"`
}

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var bottles []Bottle

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func returnAllBottles(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnAllArticles")
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(bottles)
}

func returnBottleById(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Endpoint Hit: returnBottleById")
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)

    for _, item := range bottles {
        if item.Id == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
}

func handleRequests() {
    // creates a new instance of a mux router
    myRouter := mux.NewRouter().StrictSlash(true)
    // replace http.HandleFunc with myRouter.HandleFunc
    myRouter.HandleFunc("/", homePage).Methods("GET");
    myRouter.HandleFunc("/all", returnAllBottles).Methods("GET")
    myRouter.HandleFunc("/byId/{id}", returnBottleById).Methods("GET")
    // finally, instead of passing in nil, we want
    // to pass in our newly created router as the second
    // argument
    log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
    fmt.Println("Rest API v2.0 - Mux Routers")
    bottles = []Bottle{
        Bottle{Id: "0", Title: "example", Tag: "misc", Content: "hello world"},
        Bottle{Id: "1", Title: "example", Tag: "misc", Content: "hello world!"},
    }
    handleRequests()
}
