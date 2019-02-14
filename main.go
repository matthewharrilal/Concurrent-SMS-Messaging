package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"net/url"
	"net/http"
	"github.com/joho/godotenv"
	"encoding/json"
)

func AccountConfiguration() (string, string, string) {
	// Loads our environement variables and configures url that we are going to be pinging

	err := godotenv.Load() // First load environment variables file
	if err != nil {
		log.Fatal(err)
	}

	accountSID := os.Getenv("ACCOUNT_SID")
	authToken := os.Getenv("AUTH_TOKEN")

	url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%v/Messages.json", accountSID)

	return accountSID, authToken, url
}

func ConstructMessage() strings.Reader {
	// Constructs message object with given source and destination
	// _, _, url := AccountConfiguration()

	messageData := url.Values{} // Used to store and encode following parameters to be sent over the network
	destinationNumber, sourceNumber := "7183009363", "6468324582"
	messageStub := "You are receiving a test message"
	
	// Setting source number and destination number
	messageData.Set("From", sourceNumber)
	messageData.Set("To", destinationNumber)

	
	messageData.Set("Body", messageStub)

	messageDataReader := *strings.NewReader(messageData.Encode())
	fmt.Printf("Message Data Reader ", messageDataReader)
	return messageDataReader
}

func ConstructRequest() ( http.Client, *http.Request) {
	accountSID, authToken, urlString := AccountConfiguration()
	messageDataReader := ConstructMessage()

	client := http.Client{} // In charge of executing the request
	
	// Formulate POST request with the given url string, and the encoded representation of the message body
	request, _ := http.NewRequest("POST", urlString, &messageDataReader) // Passing the message data reader by reference

	// Adds header field with the key name 'Authorization' and the two credentials we send as values to the Twillio API
	request.SetBasicAuth(accountSID, authToken)

	// Additional header fields to accept json media types which can be used for the response
	request.Header.Add("Accept", "application/json")

	// To indicate the media type that is being sent through the request
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return client, request
}

// What pairing does this function return
func ExecuteRequest() {
	// Access to the request executor and the request itself with configurations already implemented
	client, request := ConstructRequest()

	response, err := client.Do(request) // Execute the request and store the response

	// If there was an error executing the request
	if err != nil {
		log.Fatal(err)
	}

	// Checking if response came back successful
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		// Data consisting of string keys and dynamic value types depending on the JSON coming back
		data := make(map[string]interface{})

		// Decode the response body
		decoder := json.NewDecoder(response.Body)
		fmt.Printf("DECODER OF THE RESPONSE BODY ", decoder)

		err := decoder.Decode(&data) // Read the decoded data into our data map

		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Printf("Decoded Data ", data)
		}

	}
}