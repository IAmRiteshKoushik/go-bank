package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// TIP : Order or arguments is very important in Go-lang, do not change them
// http.ResponseWriter, *http.Request must come in this order because a
// http.HandlerFunc type must be satisfied and that has this order of
// elements

type APIFunc func(http.ResponseWriter, *http.Request) error

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
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// Handling error for handler
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
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
    router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccount))

	log.Println("JSON api server running on PORT: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

// Primary handler - With MUX router (unlike Gin-Gonic) we cannot specify whether
// the request is a  GET, POST or DELETE. Hence, we must handle them explicitly
// with a primary function handler
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method not allowed: %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
    // Mux vars is used to handle variables that are sent as 
    // parameters/variables (not query)
    // eg: /account/{id} -> vars["id"]
    id := mux.Vars(r)["id"]
    fmt.Println(id)
    return WriteJSON(w, http.StatusOK, &Account{})
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
