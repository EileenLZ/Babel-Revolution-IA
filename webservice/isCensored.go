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

	//récupration de la session
	session, ok := rsa.Sessions[req.Session]

	if !ok {
		w.WriteHeader(http.StatusConflict)
		msg := "This session doesn't exists. Please create a session first"
		w.Write([]byte(msg))
		return
	}

	censor := session.censorship

	//ajout du message au corpus
	censor.Corpus = append(censor.Corpus, req.Message)

	// traitement de la requête
	var resp MessageResponse = MessageResponse{req.Message, req.Title, req.Author, req.Parent, req.Session, req.Side, false, false, censorship.Events{}}
	fmt.Print(req.Side)
	if req.Side == "kgachwfxlvawkay" {
		is_message_performative := censor.IsActionPerformed(req.Title)

		if is_message_performative {
			//déclencer evt
			fmt.Print("triggered")
			resp.TriggerNewEvent = true
			resp.Events = session.Scenario.Steps[session.step].Events
			session.NextStep()
		}

	}
	if !resp.TriggerNewEvent {
		is_message_censored, censored_message, err := censor.CensordMessage(req.Message)
		is_title_censored, censored_title, err1 := censor.CensordMessage(req.Title)

		if is_message_censored || is_title_censored {
			resp.IsCensored = true
			resp.Title = censored_title     //censored_title
			resp.Message = censored_message //"L'utilisateur.ice qui a posté ce message est contrevenu.e aux textes de loi en vigueur sur la pacification des moyens de communication." //censored_message

		}
		if err != nil || err1 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("An error occured : '%s'.", err.Error())
			w.Write([]byte(msg))
			return
		}

	}

	w.WriteHeader(http.StatusOK)
	serial, _ := json.Marshal(resp)
	w.Write(serial)
}
