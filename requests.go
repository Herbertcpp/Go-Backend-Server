package main 

import (
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
)

type Response struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

func main() {
	authenticateURL := "http://localhost:8080/authenticate"

	dataMapped := map[string]interface{}{
		"username" : "herbert",
		"password" : "one",
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
