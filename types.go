package main

import "math/rand/v2"

// Struct tags will effectively tell the program that if this struct 
// gets serialized into a JSON then how the field names are going to be
type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
}

func NewAccount(firstName, lastName string) *Account {
    return &Account{
        ID: rand.IntN(10_000),
        FirstName: firstName,
        LastName: lastName,
        Number: int64(rand.IntN(1000000)),
        Balance: 0,
    }
}
