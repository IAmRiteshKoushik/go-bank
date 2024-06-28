package main

import (
    "math/rand/v2"
    "time"

    "golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
    Number      int64   `json:"number"`
    Password    string  `json:"password"`
}

// Struct tags will effectively tell the program that if this struct 
// gets serialized into a JSON then how the field names are going to be
type Account struct {
    ID        int       `json:"id"`
    FirstName string    `json:"first_name"`
    LastName  string    `json:"last_name"`
    Number    int64     `json:"number"`
    EncryptedPassword string `json:"-"`
    Balance   int64     `json:"balance"`
    CreatedAt time.Time `json:"created_at"`
}

type CreateAccountRequest struct {
    // The reason why we have this type and we are not using the Account 
    // struct is because we would be setting up the ID, Number and Balance
    // in the backend and will not be receiving them from the front-end
    FirstName string `json:"first_name"`
    LastName string `json:"last_name"`
    Password string `json:"password"`
}

type TransferRequest struct {
    ToAccount   int     `json:"to_account"`
    Amount      int     `json:"amount"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
    encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    return &Account{
        FirstName: firstName,
        LastName: lastName,
        Number: int64(rand.IntN(1000000)),
        EncryptedPassword: string(encpw),
        Balance: 0,
        CreatedAt: time.Now().UTC(),
    }, nil
}

