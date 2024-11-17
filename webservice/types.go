package webservice

type MessageRequest struct {
	Message string `json:"text"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Parent  string `json:"parent"`
	Session string `json:"session"`
	Side    string `json:"side"`
}

type Request interface {
	MessageRequest
}
