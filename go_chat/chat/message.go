package main

type Message struct {
	UserName  string `json: "UserName"`
	Body      string `json:"body"`
	Timestamp string `json:"timestamp"`
}

func (self *Message) String() string {
	return self.UserName + " at " + self.Timestamp + " says " + self.Body
}
