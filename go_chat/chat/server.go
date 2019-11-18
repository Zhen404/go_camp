package chat

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Server struct {
	connectedUsers     map[int]*User
	Messages           []*Message `json: messages`
	addUser            chan *User
	removeUser         chan *User
	newIncomingMessage chan *Message
	errorChannel       chan error
	doneCh             chan bool
}

// Create a server
func NewServer() *Server {
	Messages := []*Message{}
	connectedUsers := make(map[int]*User)
	addUser := make(chan *User)
	removeUser := make(chan *User)
	newIncomingMessage := make(chan *Message)
	errorChannel := make(chan error)
	doneCh := make(chan bool)

	return &Server{
		connectedUsers,
		Messages,
		addUser,
		newIncomingMessage,
		errorChannel,
		doneCh,
	}
}

func (server *Server) AddUser(user *User) {
	log.Println("Adding User")
	server.addUser <- user
}

func (server *Server) RemoveUser(user *User) {
	log.Println("Removing User")
	server.removeUser <- user
}

func (server *Server) Done() {
	server.doneCh <- true
}

func (server *Server) sendPastMessages(user *User) {
	for _, msg := range server.Messages {
		user.Write(msg)
	}

}

func (server *Server) Err(err error) {
	server.errorChannel <- err
}

func (server *Server) sendAll(msg *Message) {
	log.Println("In Sending to all Connected users")
	for _, user := range server.connectedUsers {
		user.Write(msg)
	}
}

func (server *Server) Listen() {
	log.Println("Server Listening ...")

	http.HandleFunc("/chat", server.handleChat)
	http.HandleFunc("/getAllMessages", server.handleGetAllMessages)

	for {
		select {
		// Adding a new user
		case user := <-server.addUser:
			log.Println("Added a new User")
			server.connectedUsers[user.id] = user
			log.Println("Now ", len(server.connectedUsers), "users are connected to chat room")
			server.sendPastMessages(user)
		case user := <-server.removeUser:
			log.Println("Removing user from chat room")
			delete(server.connectedUsers, user.id)
		case msg := <-server.newIncomingMessage:
			server.Messages = append(server.Messages, msg)
			server.sendAll(msg)
		case err := <-server.errorChannel:
			log.Println("Error :", err)
		case <-server.doneCh:
			return
		}
	}

}

func (server *Server) handleChat(responseWriter http.ResponseWriter, request *http.Request) {
	log.Println("Handling chat request")
	var messageObject Message
	conn, _ := upgrader.Upgrade(requestWriter, request, nil)
	//msgType, msg, err := conn.ReadMessage()
	err := conn.ReadJSON(&messageObject)
	log.Println("Message retireved when add user recieved", &messageObject)

	if err != nil {
		log.Println("Error while reading JSON from websocket", err.Error())
	}
	user := NewUser(conn, server)

	log.Println("going to add user", user)
	server.AddUser(user)

	log.Println("user added successfully")
	server.ProcessNewIncomingMessage(&messageObject)
}

func (server *Server) handleGetAllMessages(responseWriter http.ResponseWriter, request *http.Request) {
	json.NewEncoder(responseWriter).Encode(server)
}
