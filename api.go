package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


type APIFunc func(http.ResponseWriter, *httpRequesst) error

type APIServer struct {
    listenAddr string
}

type APIError struct {
    Error string
}


// Creating resopnse JSONs
func WriteJSON(w http.ResponseWriter, status int, v any) error {
    w.WriteHeader(status)
    w.Header().Set("Content-Type", "application/json")
    return json.NewEncoder(w).Encode(v)
}

// Making an HTTP handler as our current handlers(controllers) return error
// and that is not the defined type for an HTTP-handler according to Mux
func makeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
    // Function is an argument in the previous function, when this happens
    // we do not pass arguments
    return func(r *http.Request, w http.ResponseWriter){
        if err := f(w, r); err != nil {
            WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error())}
        }
    }
}

// Creating new API servers
func NewAPIServer(listenAddr string) *APIServer {
    return &APIServer{
        listenAddr,
    }
}

// Server initiator
func (s *APIServer) Run() {
    router := mux.NewRouter()
    router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
    log.Println("JSON api server running on PORT: ", s.listenAddr)
    http.ListenAndServe(s.listenAddr, router)
}

// Primary handler
func (s *APIServer) handleAccount(r *http.Request, w http.ResponseWriter) error {
    if r.Method == "GET" {
        return s.handleGetAccount(r, w)
    } 
    if r.Method == "POST" {
        return s.handleCreateAccount(r, w)
    }
    if r.Method == "DELETE" {
        return s.handleDeleteAccount(r, w)
    }
    return fmt.Errorf("Method not allowed: %s", r.Method)
}

func (s *APIServer) handleGetAccount(r *http.Request, w http.ResponseWriter) error {
    return nil
}

func (s *APIServer) handleCreateAccount(r *http.Request, w http.ResponseWriter) error {
    return nil
}

func (s *APIServer) handleDeleteAccount(r *http.Request, w http.ResponseWriter) error {
    return nil
}

func (s *APIServer) handleTransfer(r *http.Request, w http.ResponseWriter) error {
    return nil
}
