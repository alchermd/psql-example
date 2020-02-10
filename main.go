// An example of using the github.com/lib/pq package to
// interface with a PostgreSQL database
package main

import "log"
import "time"
import "database/sql"
import _ "github.com/lib/pq"

// Create a new data type to represent a "User" within our system
type user struct {
	id        int
	username  string
	password  string
	createdAt time.Time
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

	john := user{
		username:  "johdoe",
		password:  "secret",
		createdAt: time.Now(),
	}

	_, err = CreateUser(db, john)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(john)

	users, err := GetAllUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(users)
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
func CreateUser(db *sql.DB, u user) (sql.Result, error) {
	query := `INSERT INTO users(username, password, created_at) VALUES($1, $2, $3);`
	return db.Exec(query, u.username, u.password, u.createdAt.Format("01-02-2006"))
}

// Fetches all users and returns then as an array
func GetAllUsers(db *sql.DB) ([]user, error) {
	rows, err := db.Query(`SELECT * FROM users`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var u user
		err = rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	err = rows.Err()

	return users, err
}
