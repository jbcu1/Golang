package models

//RestResponse is json answer
type RestResponse struct {
	Result []Update `json:"result"`
}

//Update is json struct
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

//Message is json struct
type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

//Chat is json struct
type Chat struct {
	ChatID int `json:"id"`
}

//BotMessage is tratata
type BotMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}
