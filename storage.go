package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// The "lib/pq" package is imported but escaped because we do not need anything
// from the package itself, we just need it to initialize the PostgreSQL
// driver which will interact with the "database/sql" package.

type Storage interface {
    CreateAccount(*Account) error
    DeleteAccount(int) error
    UpdateAccount(*Account) error
    GetAccounts()([]*Account, error)
    GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
    db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
    connStr := os.Getenv("DATABASE_URL")
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

func (s *PostgresStore) Init() error {
    return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
    query := `CREATE TABLE IF NOT EXISTS account(
        id SERIAL PRIMARY KEY,
        first_name VARCHAR(50),
        last_name VARCHAR(50),
        number SERIAL,
        balance SERIAL,
        created_at TIMESTAMP
    )`
    _, err := s.db.Exec(query)
    return err
}


// CRUD operations
func (s *PostgresStore) CreateAccount(acc *Account) error {
    query := `
    INSERT INTO account 
    (first_name, last_name, number, balance, created_at)
    VALUES 
    ($1, $2, $3, $4, $5)`
    resp, err := s.db.Query(
        query, 
        acc.FirstName, 
        acc.LastName, 
        acc.Number, 
        acc.Balance, 
        acc.CreatedAt)
    if err != nil {
        return err
    }
    fmt.Printf("%+v\n", resp)
        
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

func (s *PostgresStore) GetAccount() ([]*Account, error) {
    rows, err := s.db.Query("SELECT * FROM account")
    if err != nil {
        return nil, err
    }

    accounts := []*Account{}
    for rows.Next() {
        account := new(Account) 
        err := rows.Scan(
            &account.ID,
            &account.FirstName,
            &account.LastName,
            &account.Number,
            &account.Balance,
            &account.CreatedAt) 
            if err != nil {
                return nil, err
            }

            accounts = append(accounts, account)
        }
    return accounts, nil
}
