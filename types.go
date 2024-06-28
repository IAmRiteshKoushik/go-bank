package main

import (
    "math/rand/v2"
    "time"
)

// Struct tags will effectively tell the program that if this struct 
// gets serialized into a JSON then how the field names are going to be
type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
    CreatedAt time.Time `json:"created_at"`
}

type CreateAccountRequest struct {
    // The reason why we have this type and we are not using the Account 
    // struct is because we would be setting up the ID, Number and Balance
    // in the backend and will not be receiving them from the front-end
    FirstName string `json:"first_name"`
    LastName string `json:"last_name"`
}

type TransferRequest struct {
    ToAccount   int     `json:"to_account"`
    Amount      int     `json:"amount"`
}

func NewAccount(firstName, lastName string) *Account {
    return &Account{
        FirstName: firstName,
        LastName: lastName,
        Number: int64(rand.IntN(1000000)),
        Balance: 0,
        CreatedAt: time.Now().UTC(),
    }
}

