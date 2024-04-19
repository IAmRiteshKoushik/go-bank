package main

import(
    "database/sql"
    _ "github.com/lib/pq"
)

// The "lib/pq" package is imported but escaped because we do not need anything
// from the package itself, we just need it to initialize the PostgreSQL 
// driver which will interact with the "database/sql" package.

type Storage interface {
    CreateAccount(*Account) error
    DeleteAccount(int) error
    UpdateAccount(*Account) error
    GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
    db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
    connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"    
    db, err := sql.Open("postgres", connStr)
    // Check for error during connection
    if err != nil {
        return nil, err
    }
    // Checking for error post connection
    if err := db.Ping(); err != nil {
        return nil, err
    }

    // Return
    return &PostgresStore{
        db: db,
    }, nil
}

func (s *PostgresStore) CreateAccount(*Account) error {
    return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
    return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
    return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
    return nil, nil
}
