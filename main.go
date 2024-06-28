package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, fname, lname, pw string) (*Account) {
    acc, err := NewAccount(fname, lname, pw)
    if err != nil {
        log.Fatal(err)
    }
    if err := store.CreateAccount(acc); err != nil {
        log.Fatal(err)
    }

    return acc
}

func seedAccounts(s Storage) {
    seedAccount(s, "Ritesh", "Koushik", "hello123")
}

func main() {
    // This allows you to create command line flags just like CLI apps
    seed := flag.Bool("seed", false, "seed the DB")
    flag.Parse()

    store, err := NewPostgresStore()
    if err != nil {
        log.Fatal(err)
    }
    if err := store.Init(); err != nil {
        log.Fatal(err)
    }

	// seed stuff
    if *seed {
        fmt.Println("Seeding the database")
        seedAccounts(store)
    }

	server := NewAPIServer(":3000", store)
	server.Run()
}
