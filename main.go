// An example of using the github.com/lib/pq package to
// interface with a PostgreSQL database
package main

import "log"
import "time"
import "net/http"
import "database/sql"
import _ "github.com/lib/pq"
import "github.com/gorilla/mux"
import "html/template"

// Create a new data type to represent a "User" within our system
type User struct {
	Id        int
	Username  string
	Password  string
	CreatedAt time.Time
}

type UsersIndexData struct {
	Users []User
}

func main() {
	// This program assumes that a database named 'example' is already created.
	// Replace the ??? part in the connection string below with your proper database credentials
	connString := `
        user=???
        password=??? 
        dbname=example 
        sslmode=disable 
        port=5432 
        host=/var/run/postgresql
    `
	db, err := GetDb(connString)
	if err != nil {
		log.Fatal(err)
	}

	_, err = CreateUsersTable(db)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", HomePageHandler).Methods("GET")
	r.HandleFunc("/users", UsersIndexHandler(db)).Methods("GET")
	r.HandleFunc("/users", StoreUserHandler(db)).Methods("POST")
	r.HandleFunc("/users/new", CreateUserHandler).Methods("GET")

	log.Print("Starting server")
	http.ListenAndServe(":8000", r)
}

// HTTP handler for the home page
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	homePageView := template.Must(template.ParseFiles("views/index.html"))
	homePageView.Execute(w, nil)
}

// Handler for storing new users.
// Wraps the expect function signature in order to accept a database instance
func StoreUserHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := User{
			Username:  r.FormValue("username"),
			Password:  r.FormValue("password"),
			CreatedAt: time.Now(),
		}

		_, err := CreateUser(db, u)

		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/users", 302)
	}
}

// Handler for displaying all the users.
// Wraps the expect function signature in order to accept a database instance
func UsersIndexHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := GetAllUsers(db)
		if err != nil {
			log.Fatal(err)
		}

		usersIndexView := template.Must(template.ParseFiles("views/users/index.html"))
		usersIndexView.Execute(w, UsersIndexData{
			Users: users,
		})
	}
}

// Handler for showing the form for creating a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	createUserView := template.Must(template.ParseFiles("views/users/create.html"))
	createUserView.Execute(w, nil)
}

// Opens a postgresql connection and returns a pointer to a new database instance
func GetDb(connString string) (*sql.DB, error) {
	return sql.Open("postgres", connString)
}

// Creates a users table if it doesn't exist yet
func CreateUsersTable(db *sql.DB) (sql.Result, error) {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT NOT NULL,
            password TEXT NOT NULL,
            created_at TIMESTAMP
        );
    `
	return db.Exec(query)
}

// Inserts a new row using the information from the provided user instance
func CreateUser(db *sql.DB, u User) (sql.Result, error) {
	query := `INSERT INTO users(username, password, created_at) VALUES($1, $2, $3);`
	return db.Exec(query, u.Username, u.Password, u.CreatedAt.Format("01-02-2006 15:04:05"))
}

// Fetches all users and returns then as an array
func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`SELECT * FROM users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err = rows.Scan(&u.Id, &u.Username, &u.Password, &u.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	err = rows.Err()

	return users, err
}
