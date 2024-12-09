package webservice

import "TestNLP/pkg/censorship"

type Side int64

const (
	QG      Side = 0
	Terrain Side = 1
)

type MessageRequest struct {
	Message string `json:"text"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Parent  string `json:"parent"`
	Session string `json:"session"`
	Side    string `json:"side"`
}

type MessageResponse struct {
	Message         string            `json:"text"`
	Title           string            `json:"title"`
	Author          string            `json:"author"`
	Parent          string            `json:"parent"`
	Session         string            `json:"session"`
	Side            string            `json:"side"`
	IsCensored      bool              `json:"isCensored"`
	TriggerNewEvent bool              `json:"triggerNewEvent"`
	Events          censorship.Events `json:"events"`
}

type Request interface {
	MessageRequest | NewSessionRequest
}

type NewSessionRequest struct {
	Session     string   `json:"session"`
	BannedWords []string `json:"bannedWords"`
}

type NewSessionResponse struct {
	Session string `json:"session"`
}
