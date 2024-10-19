package webservice

type MessageRequest struct {
	Message string `json:"message"`
}

type MessageResponse struct {
	IsCensored bool `json:"is_censored"`
}

type Request interface {
	MessageRequest
}
