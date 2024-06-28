package main

import (
	"database/sql"
	"fmt"
	"os" 

    "github.com/joho/godotenv"
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
    GetAccountByNumber(int) (*Account, error)
}

type PostgresStore struct {
    db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
    // Incase environment variables do not load properly, then handle the error
    err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file", err)
        return nil, err
    }
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
    // for initializing a database, there should be a default accounts table 
    // ready to accept incoming data.
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
    _, err := s.db.Query(
        query, 
        acc.FirstName, 
        acc.LastName, 
        acc.Number, 
        acc.Balance, 
        acc.CreatedAt)
    if err != nil {
        return err
    }
    return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
    return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
    // After deleting a field, you need not return the deleted field but just 
    // the confirmation of whether they have been deleted or not.
    _, err := s.db.Query("DELETE FROM account WHERE id = $1", id)
    return err
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
    rows, err := s.db.Query("SELECT * FROM account WHERE id = $1", id)
    if err != nil {
        return nil, err
    }
    for rows.Next(){
        // scanIntoAccount directly returns an account pointer and error
        return scanIntoAccount(rows)
    }
    // if there was no rows.Next() then the table did not contain 
    // a single row which matched the particular ID, in which case,
    // we do not need to return any pointer but we must return an error
    return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
    rows, err := s.db.Query("SELECT * FROM account WHERE number = $1", number)
    if err != nil {
        return nil, err
    }

    for rows.Next(){
        return scanIntoAccount(rows)
    }
    return nil, fmt.Errorf("Account with number [%d] not found", number)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
    // Fetching all rows from the account table 
    rows, err := s.db.Query("SELECT * FROM account")
    if err != nil {
        return nil, err
    }

    // After fetching everything from the account table,
    // we need to move everything to the slice of account 
    // pointers. 
    accounts := []*Account{}
    for rows.Next() {
        account, err := scanIntoAccount(rows)
        if err != nil {
            return nil, err
        }
        accounts = append(accounts, account)
    }
    return accounts, nil
}

// -- HELPER FUNCTION useful for getting things from SQL rows 
// and moving them into Account struct and returning a pointer.
// Will be useful in other functions as well.
func scanIntoAccount(rows *sql.Rows) (*Account, error){
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
    return account, nil
}
