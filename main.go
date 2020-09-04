package main

import (
	"RESTGo/tokens"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

var(
	MessagesData []MessageData
	PageMessagesData = []PageMessageData{
		PageMessageData{Username: "SYSTEM", Message:"Welcome to the Gamer zone!"},
	}

	messageNumber = 0
)


//Convert Tokens to usernames
func tokToUsername(t string, writer http.ResponseWriter) string {
	if t == tokens.WillToken {
		return "Will"
	}
	if t == tokens.JakeToken {
		return "Jake"
	}
	if t == tokens.DawsonToken {
		return "Dawson"
	} else {
		writer.WriteHeader(418)
		return ""
	}

}

type MessageData struct {
	UserID string `json:"UserID"`
	Message string `json:"Message"`
}

type PageMessageData struct {
	Username string `json:"Username"`
	Message string `json:"Message"`
}

func defaultPage(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Go to /board")
	fmt.Println("Request received: default page")

}

func allMessages(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Request received: All Messages")
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(PageMessagesData)
	if err != nil {
		fmt.Println(err)
	}
}

func createMessage(writer http.ResponseWriter, request *http.Request) {
	reqBody, _ := ioutil.ReadAll(request.Body)
	var message MessageData
	errUnmarsh := json.Unmarshal(reqBody, &message)
	if errUnmarsh != nil {
		fmt.Println(errUnmarsh)
		writer.WriteHeader(400)
		return
	}

	if !(len(message.Message) > 0)  {
		writer.WriteHeader(400)
		return
	}

	//fmt.Println(message)
	username := tokToUsername(message.UserID, writer)
	if username == "" {
		return
	}

	var newMessage = make([]PageMessageData,1)
	newMessage[0].Username = username
	newMessage[0].Message = message.Message

	PageMessagesData = append(PageMessagesData, newMessage[0])
	messageNumber++
	errEncode := json.NewEncoder(writer).Encode(message)
	if errEncode != nil {
		fmt.Println(errEncode)
		writer.WriteHeader(500)
		return
	}
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", defaultPage)
	router.HandleFunc("/board", allMessages)
	router.HandleFunc("/send", createMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":1338", router))
}

func main() {
	handleRequests()
}