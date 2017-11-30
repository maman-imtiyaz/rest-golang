package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"./lib"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	ID        string
	Firstname string
	Lastname  string
	Address   *Address
}
type Address struct {
	City  string
	State string
}

type Result struct {
	Status bool
	Result []Person
	Message lib.Exception
}

var result = Result{}
var people []Person

func GetPeople(w http.ResponseWriter, r *http.Request) {
	lib.Block{
		Try: func() {
			result.Result = people
			result.Status = true
			result.Message = "Success"
		},
		Catch: func(e lib.Exception) {
			result.Status = false
			result.Message = e
		},
		Finally: func() {
			json.NewEncoder(w).Encode(result)
		},
	}.Do()
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var temp []Person
	lib.Block{
		Try: func() {
			for _, item := range people {
				if item.ID == params["id"] {
					temp = append(temp, item)
					result.Result = temp
					return
				}
			}
			result.Status = true
			result.Message = "Success"
		},
		Catch: func(e lib.Exception) {
			result.Status = false
			result.Message = e
		},
		Finally: func() {
			json.NewEncoder(w).Encode(result)
		},
	}.Do()
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	lib.Block{
		Try: func() {
			var person Person
			person.ID = r.FormValue("id")
			person.Firstname = r.FormValue("firstname")
			person.Lastname = r.FormValue("lastname")
			people = append(people, person)

			result.Result = people
			result.Status = true
			result.Message = "Success"
		},
		Catch: func(e lib.Exception) {
			result.Status = false
			result.Message = e
		},
		Finally: func() {
			json.NewEncoder(w).Encode(result)
		},
	}.Do()
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lib.Block{
		Try: func() {
			for index, item := range people {
				if item.ID == params["id"] {
					people = append(people[:index], people[index+1:]...)
					break
				}
			}
			result.Result = people
			result.Status = true
			result.Message = fmt.Sprintf("Success deleted id %s", params["id"])
		},
		Catch: func(e lib.Exception) {
			result.Status = false
			result.Message = e
		},
		Finally: func() {
			json.NewEncoder(w).Encode(result)
		},
	}.Do()
}

func MgoGetPeople(w http.ResponseWriter, r *http.Request) {
	lib.Block{
		Try: func() {
			temp := []Person{}
			session, _ := mgo.Dial("localhost")
			c := session.DB("db_training").C("training")
			iter := c.Find(bson.M{}).Iter()

			person := Person{}
			for iter.Next(&person){
				temp = append(temp, person)
			}

			session.Close()
			result.Result = temp
			result.Status = true
			result.Message = "Success"
		},
		Catch: func(e lib.Exception) {
			result.Status = false
			result.Message = e
		},
		Finally: func() {
			json.NewEncoder(w).Encode(result)
		},
	}.Do()
}

func MgoGetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lib.Block{
		Try: func() {
			temp := []Person{}
			session, err := mgo.Dial("localhost")
			c := session.DB("db_training").C("training")
			iter := c.Find(bson.M{"id": params["id"]}).Iter()

			if err != nil{
				panic(err)
			}

			person := Person{}
			for iter.Next(&person){
				temp = append(temp, person)
			}

			session.Close()
			result.Result = temp
			result.Status = true
			result.Message = "Success"
		},
		Catch: func(e lib.Exception) {
			result.Status = false
			result.Message = e
		},
		Finally: func() {
			json.NewEncoder(w).Encode(result)
		},
	}.Do()
}

func MgoCreatePerson(w http.ResponseWriter, r *http.Request) {
	lib.Block{
		Try: func() {
			person := Person{}
			temp := []Person{}
			//var temp []Person
			person.ID = r.FormValue("id")
			person.Firstname = r.FormValue("firstname")
			person.Lastname = r.FormValue("lastname")

			session, err := mgo.Dial("localhost")
			c := session.DB("db_training").C("training")
			c.Insert(person)
			session.Close()

			if err != nil{
				panic(err)
			}

			temp = append(temp, person)
			result.Result = temp
			result.Status = true
			result.Message = fmt.Sprintf("%s Success inserted to mongo.", person.Firstname)
		},
		Catch: func(e lib.Exception) {
			result.Status = false
			result.Message = e
		},
		Finally: func() {
			json.NewEncoder(w).Encode(result)
		},
	}.Do()
}

func MgoDeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	lib.Block{
		Try: func() {
			temp := []Person{}
			session, err := mgo.Dial("localhost")
			c := session.DB("db_training").C("training")
			c.Remove(bson.M{"id": params["id"]})
			iter := c.Find(nil).Iter()
			session.Close()

			if err != nil{
				panic(err)
			}

			person := Person{}
			for iter.Next(&person){
				temp = append(temp, person)
			}

			result.Result = temp
			result.Status = true
			result.Message = fmt.Sprintf("Success deleted id %s", params["id"])
		},
		Catch: func(e lib.Exception) {
			result.Status = false
			result.Message = e
		},
		Finally: func() {
			json.NewEncoder(w).Encode(result)
		},
	}.Do()
}

// our main function
func main() {
	people = append(people, Person{ID: "1", Firstname: "Maman", Lastname: "Lesmana", Address: &Address{City: "Bandung", State: "Indonesia"}})
	people = append(people, Person{ID: "2", Firstname: "Slamet", Lastname: "Rohadi", Address: &Address{City: "Jambi", State: "Indonesia"}})
	people = append(people, Person{ID: "3", Firstname: "Vriandy", Lastname: "Raplh"})

	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")

	router.HandleFunc("/mgo_people", MgoGetPeople).Methods("GET")
	router.HandleFunc("/mgo_people/{id}", MgoGetPerson).Methods("GET")
	router.HandleFunc("/mgo_people", MgoCreatePerson).Methods("POST")
	router.HandleFunc("/mgo_people/{id}", MgoDeletePerson).Methods("DELETE")

	fmt.Println("Running at localhost port 8000..")
	log.Fatal(http.ListenAndServe(":8000", router))
}
