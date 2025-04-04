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
	id       string
	addr     string
	Sessions map[string]Session
}

func NewServerAgent(addr string) *ServerAgent {
	return &ServerAgent{sync.Mutex{}, addr, addr, map[string]Session{}}
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
	mux.HandleFunc("POST /api/checkMsg", rsa.DoIsCensored)
	mux.HandleFunc("POST /api/newSession", rsa.DoNewSession)
	mux.HandleFunc("GET /api/health", rsa.Health)

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

func (rsa *ServerAgent) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
