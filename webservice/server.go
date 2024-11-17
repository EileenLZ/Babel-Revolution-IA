package webservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type ServerAgent struct {
	sync.Mutex
	id          string
	addr        string
	bannedWords []string
}

func NewServerAgent(addr string) *ServerAgent {
	return &ServerAgent{id: addr, addr: addr, bannedWords: []string{
		"clavier", "parapluie", "flaque", "écran",
		"machine", "IA", "SOPHIA", "détruire",
	}}
}

// Test de la méthode
func (rsa *ServerAgent) checkMethod(method string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "method %q not allowed", r.Method)
		return false
	}
	return true
}

// do a decoderequest factory
func decodeRequest[Req Request](r *http.Request) (req Req, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	err = json.Unmarshal(buf.Bytes(), &req)
	return
}

func (rsa *ServerAgent) Start() {
	// création du multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/CheckMsg", rsa.DoIsCensored)

	// création du serveur http
	s := &http.Server{
		Addr:           rsa.addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	// lancement du serveur
	log.Println("Listening on", rsa.addr)
	go log.Fatal(s.ListenAndServe())
}
