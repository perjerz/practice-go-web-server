package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Student struct {
	Name  string  `json:"name"`
	Age   int     `json:"age"`
	Score float64 `json:"score"`
}

type Todo struct {
	UserId    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	r := mux.NewRouter()
	port := 8080
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/decode", decode).Methods("POST")
	r.HandleFunc("/unmarshal", unmarshal).Methods("POST")
	r.HandleFunc("/encode", encode).Methods("GET")
	r.HandleFunc("/marshal", marshal).Methods("GET")
	r.HandleFunc("/fetch", fetchGoogle).Methods("GET")
	r.HandleFunc("/student/{student}/todo/{todo}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		s := vars["student"]
		t := vars["todo"]
		fmt.Fprintf(w, "student: %v, todo %v ", s, t)
	}).Methods("GET")
	r.Handle("/notfound", r.NotFoundHandler)
	r.Handle("/static", http.StripPrefix("/static/", http.FileServer(http.Dir("/static/"))))

	log.Printf("Listen on port: %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello OK"))
}

func decode(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		{
			var s Student
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&s)
			if err != nil {
				http.Error(w, "Error", http.StatusBadRequest)
				return
			}
			fmt.Println(s)
			w.Write([]byte("Good job"))
		}
	default:
		w.Write([]byte("Method Not Allowed"))
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func unmarshal(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		{
			var s Student
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = json.Unmarshal(b, &s)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Println(s)
			w.Write([]byte("Good job"))
		}
	default:
		w.Write([]byte("Method Not Allowed"))
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func encode(w http.ResponseWriter, r *http.Request) {
	s := Student{
		Name:  "Srikote",
		Age:   24,
		Score: 100.0,
	}
	encoder := json.NewEncoder(w)
	err := encoder.Encode(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func marshal(w http.ResponseWriter, r *http.Request) {
	s := Student{
		Name:  "Supanut",
		Age:   23,
		Score: 1000.0,
	}
	encoder := json.NewEncoder(w)
	err := encoder.Encode(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func fetchGoogle(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/1")
	//resp, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/todos/1", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err.Error())
		return
	}

	var t Todo
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&t)

	encoder := json.NewEncoder(w)
	encoder.Encode(t)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err.Error())
		return
	}
}
