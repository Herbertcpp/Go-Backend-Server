package main 

import (
	"fmt"
	"strings"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"log"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	db, err := sql.Open("sqlite3", "./DataBase.db")
	if err != nil {
		log.Fatal(err)
	}
	initDataBase(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/home", CORSMiddleWare(home))
	mux.HandleFunc("/register", CORSMiddleWare(func(w http.ResponseWriter, r *http.Request) {
		addUser(db, w, r)
	}))

	mux.HandleFunc("/print", CORSMiddleWare(func(w http.ResponseWriter, r *http.Request) {
		printDataBase(db, w, r)
	}))

	mux.HandleFunc("/authenticate", CORSMiddleWare(func(w http.ResponseWriter, r *http.Request) {
		authenticateUser(db, w, r)
	}))

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", mux)
}

//Handler Functions and MiddleWare

func CORSMiddleWare(next http.HandlerFunc) http.HandlerFunc {
		return(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "*")

			if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
			}
			next(w, r)
		})

		
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to server\n"))	
}

func addUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {

		w.Header().Set("Content-Type", "application/json")

		jsonErrResponse := map[string]interface{}{
			"success" : false,
			"message" : "error reading json",
		}

		encodeErr := json.NewEncoder(w).Encode(&jsonErrResponse)
		if encodeErr != nil {
			fmt.Println("Even more errors :( (Encoding json)")
		}
	}
	success := saveToDataBase(db, newUser.Username, newUser.Password)
	if success {
		fmt.Println("Successfully added", newUser.Username, "to the database")

		successJson := map[string]interface{}{
			"success" : true,
			"message" : "user added to database successfully",
		}
		encodeErrTwo := json.NewEncoder(w).Encode(successJson)
		if encodeErrTwo != nil {
			fmt.Println("Problems responding to request")
		}
	} else {
			errJson := map[string]interface{}{
			"success" : false,
			"message" : "username taken already",
		}
		encodeErrTwo := json.NewEncoder(w).Encode(errJson)
		if encodeErrTwo != nil {
			fmt.Println("Problems responding to request")
		}
	}
}

func authenticateUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var user User

	decodeErr := json.NewDecoder(r.Body).Decode(&user)
	if decodeErr != nil {
		fmt.Println("Error decoding json")
	}

	var hashedPassword string
	err := db.QueryRow("select password from users where username = ?", user.Username).Scan(&hashedPassword)

	if err == sql.ErrNoRows {
		fmt.Println("No rows found")
		noRowsJson := map[string]interface{}{
			"success" : false,
			"message" : "no user found",
		}
		json.NewEncoder(w).Encode(noRowsJson)
		return
	}

	if err != nil {
		fmt.Println("Error when reading the rows")
	}

	if !checkHashed(hashedPassword, user.Password) {
		errJson := map[string]interface{}{
			"success" : false,
			"message" : "wrong password",
		}	
		json.NewEncoder(w).Encode(errJson)
	} else {
			succ := map[string]interface{}{
			"success" : true,
			"message" : "correct password",
		}	
		json.NewEncoder(w).Encode(succ)		
	}
	
}

//Data Base Functions

func initDataBase(db *sql.DB) {
	_, err := db.Exec("create table if not exists users (id integer primary key, username text unique, password text)")	
	if err != nil {
		log.Fatal(err)
	}
}

func saveToDataBase(db *sql.DB, username string, password string) bool {
	if strings.TrimSpace(username) == "" {
		return false
	}
	hashed, hashErr := bcrypt.GenerateFromPassword([]byte(password), 10)
	if hashErr != nil {
		fmt.Println("Error hashing the password")
		return false
	}
	_, err := db.Exec("insert into users (username, password) values (?, ?)", username, hashed)
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {

	fmt.Println("Username is taken already!")

		return false
	}
	if err != nil {
		return false
	}
	return true
}

func printDataBase(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("select username, password from users")
	if err != nil {
		w.Write([]byte("Error reading users"))
	}	

	var username, password string
	var mapped = make(map[string]string)

	for rows.Next() {
		ScanErr := rows.Scan(&username, &password)
		if ScanErr != nil {
			fmt.Println("Error reading line from database")
		}
		mapped[username] = password
	}
	newEncoder := json.NewEncoder(w)
	newEncoder.SetIndent("", " ")
	sendingErr := newEncoder.Encode(mapped)
	if sendingErr != nil {
		fmt.Println("Issue sending data to the client!")
	}
}

// Bcrypt functions 

func checkHashed(hashed string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
