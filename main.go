package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "postgresql"
// 	dbname   = "bird_encyclopedia"
// )

func newRouter() *mux.Router {
	r := mux.NewRouter()

	// Static Directory for assets
	staticFileDirectory := http.Dir("./assets/")
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/assets").Handler(staticFileHandler).Methods("GET")

	// Add handlers to the router
	r.HandleFunc("/hello", handler).Methods("GET")
	r.HandleFunc("/bird", getBirdHandler).Methods("GET")
	r.HandleFunc("/bird", createBirdHandler).Methods("POST")
	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	u, err := url.Parse(dbUrl)
	if err != nil {
		panic(err)
	}

	user := u.User.Username()
	password, _ := u.User.Password()
	host, port, _ := net.SplitHostPort(u.Host)
	dbname := strings.TrimPrefix(u.Path, "/")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open(u.Scheme, psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	InitStore(&dbStore{db: db})

	r := newRouter()

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
