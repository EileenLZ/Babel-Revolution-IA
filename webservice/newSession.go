package webservice

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (rsa *ServerAgent) DoNewSession(w http.ResponseWriter, r *http.Request) {
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := decodeRequest[NewSessionRequest](r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	_, exists := rsa.Sessions[req.Session]

	if exists {
		w.WriteHeader(http.StatusBadRequest)
		msg := "This session already exists."
		w.Write([]byte(msg))
		return
	}

	new_session, err := NewSession(req.Session, scenarioFile)
	rsa.Sessions[req.Session] = *new_session

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	var resp NewSessionResponse
	resp.Session = req.Session

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}
