package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
    "os"

	"github.com/gorilla/mux"
    jwt "github.com/golang-jwt/jwt/v4"
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

    router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
    router.HandleFunc("/account/{id}", withJWT(makeHTTPHandleFunc(s.handleGetAccountByID), s.store))

    // Here, you can do "/transfer/{accountNumber}" but then if anyone checks 
    // the browser history they would be able to find the account number to 
    // which a transfer has been made. On the contrary if we do not specify that 
    // and manage it as a POST request instead of a GET request, then we would 
    // have to inspect the Network tab when this particular request is going in 
    // order to know. And this information gets deleted and not stored in the
    // browser cache.
    router.HandleFunc("/transfer", withJWT(makeHTTPHandleFunc(s.handleTransfer), s.store))

    // NOTE : AccountNumbers are safe and not hackable but that being said, in 
    // order to ensure better privacy, it is better to not have them exposed.

	log.Println("JSON api server running on PORT", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
    if r.Method != "POST"{
        return fmt.Errorf("method not allowed %s", r.Method)
    }
    req := new(LoginRequest)
    if err := json.NewDecoder(r.Body).Decode(req); err != nil {
        return err
    }

    // search for the user 

    return WriteJSON(w, http.StatusOK, req)
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
    req := new(CreateAccountRequest)
    if err := json.NewDecoder(r.Body).Decode(req); err != nil {
        return err
    }
    account, err := NewAccount(req.FirstName, req.LastName, req.Password)
    if err != nil {
        return err
    }
    if err := s.store.CreateAccount(account); err != nil {
        return err
    }

    // -- OUTDATED
    // After the account is successfully created, create a JWT token
    // tokenString, err := createJWT(account)
    // if err != nil {
    //     return err
    // }
    // fmt.Println("JWT TOken: ", tokenString)

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
    transferReq := new(TransferRequest)
    if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
        return err
    }
    // You need to clear the previous request before waiting for a new request
    // If this is not done, there will be a resource leak.
    defer r.Body.Close()

    // FIX : We need to call in storage methods
    return WriteJSON(w,http.StatusOK, transferReq)
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

func permissionDenied(w http.ResponseWriter){
    // it is a good idea to setup some form of logging and send it over to your 
    // logging / monitoring tools - DataDog (or) Grafana etc.
    WriteJSON(w, http.StatusForbidden, APIError{ Error: "permission denied" })
}

// A decorator function which is going to sit on top of handler functions 
// and authenticate before processing requests.
func withJWT(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request){
        fmt.Println("Calling JWT Auth Middleware")
        tokenString := r.Header.Get("x-jwt-token")
        token, err := validateJWT(tokenString)
        // Validate JWT only checks if the signing method works but it does 
        // return back the token in both cases which is a struct that has a 
        // 'Valid' field. An invalid token does not generate an error
        if err != nil {
            permissionDenied(w)
            return
        }
        // Here, we need to check if the token is valid or not by accessing the 
        // field inside the token-struct. After we have done so, we can proceed 
        // and check
        if !token.Valid {
            permissionDenied(w)
            return
        }
        userID, err := getID(r)
        if err != nil {
            permissionDenied(w)
            return
        }
        account, err := s.GetAccountByID(userID)
        if err != nil {
            return
        }
        // the claims are in string-format and need to be converted to a 
        // map[string]interface{} - format before we can access it.
        claims := token.Claims.(jwt.MapClaims)
        // PROBLEM : During conversion, the AccountNumber becomes a float64 
        // which falls under the interface{} implementation. But the account 
        // number that we have retrived from the database falls under the int64
        // So, in order to check the equality, we need to convert the AccountNum 
        // to int64, but we cannot do this directly because the type of AccountNum 
        // is decided in the run-time and not during compile-time. In-order to 
        // make the conversion possible, we mut first type-assert it into float64 
        // and then make the conversion to int64. In-case the float64 type assetion
        // fails, the program will panic (POTENTIAL FAILURE POINT)
        if account.Number != int64(claims["AccountNumber"].(float64)){
            permissionDenied(w)
            return
        }

        handlerFunc(w, r)
    }
}

func validateJWT(tokenString string) (*jwt.Token, error) {
    secret := os.Getenv("JWT_SECRET")

    // anonymous function is passed inside the jwt.Parse() function. If token 
    // is parsed properly, then you return back the
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
        // If the signing method does not match with the signing method setup 
        // in the backend then an error is generated and returned
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }

        // if the signing method is valid, then we need to check if the token 
        // is actually valid or not. Basically you are checkign at two levels:
        // 1. Checking if the signing algorithm is alright (VULNERABILITY)
        // 2. Checking if the token is alright
        return []byte(secret), nil
    })
}

// -- Sample JWT String (valid)
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2NvdW50TnVtYmVyIjo0MDMxMzgsIkV4cGlyZXNBdCI6MTUwMDB9.A2tRhYooBvzCUL7fKtj87TWTu_MGe4LwYFz2ies1CAc

func createJWT(account *Account) (string, error) {
    claims := &jwt.MapClaims{
        "ExpiresAt": 15000,
        "AccountNumber": account.Number,
    }
    // Get the secret from environment variables
    secret := os.Getenv("JWT_SECRET")
    // Mention the claims and the signing method
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Return back the token after signing it with the secret
    // the secret should be a byte-slice
    return token.SignedString([]byte(secret))

}
