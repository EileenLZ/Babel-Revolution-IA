package webservice

import (
	"TestNLP/pkg/censorship"
	"encoding/json"
	"fmt"
	"net/http"
)

func (rsa *ServerAgent) DoIsCensored(w http.ResponseWriter, r *http.Request) {
	// mise à jour du nombre de requêtes
	rsa.Lock()
	defer rsa.Unlock()

	// vérification de la méthode de la requête
	if !rsa.checkMethod("POST", w, r) {
		return
	}

	// décodage de la requête
	req, err := decodeRequest[MessageRequest](r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	// traitement de la requête
	var resp MessageResponse

	resp.IsCensored, err = censorship.IsSentenceCensored(req.Message, rsa.bannedWords)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("An error occured : '%s'.", err.Error())
		w.Write([]byte(msg))
		return
	}

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}
