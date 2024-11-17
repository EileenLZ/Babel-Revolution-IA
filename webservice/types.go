package webservice

import "TestNLP/pkg/censorship"

type MessageRequest struct {
	Message string `json:"text"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Parent  string `json:"parent"`
	Session string `json:"session"`
	Side    string `json:"side"`
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

type Session struct {
	censorship censorship.Censorship
	id         string
}

func NewSession(id string, banned_words []string) *Session {
	return &Session{*censorship.NewCensorship(banned_words), id}
}
