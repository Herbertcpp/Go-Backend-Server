package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func main() {

	for {
		var userInput int
		fmt.Print("Welcome to the testing script! Following options are avaviable:\n(0) Quit\n(1) Register a user\n(2) Authenticate a user\n(3) Print all users\nChoice: ")
		fmt.Scan(&userInput)

		switch userInput {
		case 0:
			fmt.Println("Quitting... See you next time!")
			return
		case 1:
			fmt.Println("1")
		case 2:
			fmt.Println("2")
		case 3:
			fmt.Println("3")
		default:
			fmt.Println("Invalid Input")
		}
	}

}

func authenticate() {
	authenticateURL := "http://localhost:8080/authenticate"

	var inputUsername, inputUserpassword string
	fmt.Print("Username: ")
	fmt.Scan(&inputUsername)
	fmt.Print("Password: ")
	fmt.Scan(&inputUserpassword)

	dataMapped := map[string]interface{}{
		"username": inputUsername,
		"password": inputUserpassword,
	}

	jsonData, JsonErr := json.Marshal(dataMapped)
	if JsonErr != nil {
		fmt.Println("Problem converting map to json")
	}

	resp, err := http.Post(
		authenticateURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		fmt.Println("Error sending json")
		return
	}
	defer resp.Body.Close()

	var response Response
	decodingErr := json.NewDecoder(resp.Body).Decode(&response)
	if decodingErr != nil {
		fmt.Println("Issue decoding file")
		return
	}
	fmt.Println(response.Success, response.Message)
}

func register() {

}
