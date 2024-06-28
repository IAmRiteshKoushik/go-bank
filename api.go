package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// TIP : Order of arguments is very important in Go-lang, do not change them
// http.ResponseWriter, *http.Request must come in this order because a
// http.HandlerFunc type must be satisfied and that has this order of
// elements

type APIFunc func(http.ResponseWriter, *http.Request) error

type APIServer struct {
	listenAddr string
    store Storage
}

type APIError struct {
	Error string `json:"error"`
}

// Server initiator
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
    router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))

	log.Println("JSON api server running on PORT", s.listenAddr)
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
    // After handling GET and POST if there are other HTTP verbs being 
    // used then those are not to be considered (for the time being)
	return fmt.Errorf("Method not allowed: %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
    accounts, err := s.store.GetAccounts()
    if err != nil {
        return err
    }
    return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
    if r.Method == "GET" {

        // Mux vars is used to handle variables that are sent as 
        // parameters/variables (not query)
        // eg: /account/{id} -> vars["id"]

        // -- OUTDATED
        // The id grabbed from the URL vars is not in integer format, we need to 
        // convert it into integer format and then utilize it, for which we can 
        // actually run a check against whether the value is garbage or not.
        // idStr := mux.Vars(r)["id"]
        // id, err := strconv.Atoi(idStr)
        // if err != nil {
        //     return fmt.Errorf("Invalid ID given %s", idStr)
        // }

        id, err := getID(r)
        if err != nil {
            return err
        }
        // After the ID is valid, we can go ahead and run a query against the 
        // database and if the query is successful, send that value to WriteJSON
        // or else return the error generated
        account, err := s.store.GetAccountByID(id)
        if err != nil {
            return err
        }

        // The next error will come in encoding the data that has come in the form 
        // of a struct into an HTTP response (JSON) format. So in-order to do that 
        // we need to utilize the WriteJSON function. Here, if the encoder works 
        // correctly, then we will send the data back as API response/
        return WriteJSON(w, http.StatusOK, account)
    }
    // We have not setup separate pathway for DELETE and the Mux router does not 
    // handle it by default so we need to wrap up inside the function which 
    // is handling the path which is already handling the "id" parameter and 
    // pass the control over to DELETE method
    if r.Method == "DELETE" {
        return s.handleDeleteAccount(w, r)
    }

    return fmt.Errorf("Method not allowed: %s", r.Method)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
    // Using the new keyword so that we get a reference to the structure and 
    // not the actual structure. (Reduces memory overhead). Also, the Decode
    // method takes in a pointer to a structure
    createAccountReq := new(CreateAccountRequest)
    if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
        return err
    }
    account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
    if err := s.store.CreateAccount(account); err != nil {
        return err
    }

    // FIX ME : Currently the database does not tell you whether the account 
    // is successfully created or not as there is only creation and not 
    // fetching of data. This will most probably require a fix in the future

    return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
    id, err := getID(r)
    if err != nil {
        return err
    }
    // if the ID is valid, then we run a check against the database and see if 
    // any error is generated or not
    if err := s.store.DeleteAccount(id); err != nil {
        return err
    }
    // If no errors are generated then we WriteJSON() or else there are top 
    // level functions which can handle the error and send back a BadRequest 
    // with the error code.
    return WriteJSON(w, http.StatusOK, map[string]int{ "deleted" : id })
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// -- Helper Functions
// Creating resopnse JSONs
func WriteJSON(w http.ResponseWriter, status int, v any) error {
    // The content-type has to come before setting adding additional info
    // as it defines the structure of the data.

    // If we do the opposite, there is "silent failure" as the header 
    // is considered committed as "text/html" type by default
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
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

// Creating new API server
func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr,
        store,
	}
}

// Getting ID from a request-URL and then converting it appropriately
func getID(r *http.Request) (int, error) {
    idStr := mux.Vars(r)["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return id, fmt.Errorf("Invalid id given %s", idStr)
    }
    return id, nil
}
