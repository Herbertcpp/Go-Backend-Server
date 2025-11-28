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
			register()
		case 2:
			authenticate()
		case 3:
			printUserData()
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
	registerURL := "http://localhost:8080/register"

	var InputUsername, Inputuserpassword string 
	fmt.Print("Username: ")
	fmt.Scan(&InputUsername)
	fmt.Print("Password: ")
	fmt.Scan(&Inputuserpassword)

	dataMapped := map[string]string{
		"username" : InputUsername,
		"password" : Inputuserpassword,
	}

	jsonData, conversionErr := json.Marshal(dataMapped)
	if conversionErr != nil {
		fmt.Println("Error converting data to json")
		return
	}
	var resp Response
	requestResponse, err := http.Post(registerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request to server :(")
	}
	decodeErr := json.NewDecoder(requestResponse.Body).Decode(&resp)
	if decodeErr != nil {
		fmt.Println("Error while decoding")
		return
	}
	if resp.Success {
		fmt.Println("Registered User: ", InputUsername)
	} else {
		fmt.Println("Username taken already")
	}
}

func printUserData() {
	url := "http://localhost:8080/print"
	var recievedData map[string]string

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending request")
	}

	decodeErr := json.NewDecoder(resp.Body).Decode(&recievedData)
	if decodeErr != nil {
		fmt.Println("Error decoding server response")
	}
	defer resp.Body.Close()

	for user, password := range recievedData {
		fmt.Printf("%s, %s\n", user, password)
	}
}
