# Overview
Simple HTTP-REST API written in Go using Gorilla/Mux, PostgreSQL and JWT Auth.

## Setup
Clone the project and run the corresponding commands
```bash
# Cloning
git clone https://github.com/IAmRiteshKoushik/go-bank

# Fetching all required dependencies
go mod tidy
go mod vendor

# Compiling and the project
make build

# Seeding with sample data : use the --seed flag
# Seed data : { FName: Ritesh, LName: Koushik, Password: hello123 }
./bin/go-bank --seed

# Running the project
make run #(or)
./bin/go-bank
```
All further testing can be run through Postman, cURL, ThunderClient etc.

## Run Locally
Create the environment variables file
```bash
touch .env
```
Setup the environment variables
```bash 
# Inside the .env file, have the following KV pairs
DATABASE_URL="<your-database-connection-string-here>"
JWT_SECRET="<your-jwt-secret-here>"
```
Have the database running. You can either have your local installation of 
PostgreSQL, a cloud provider like Neon or a docker container. It is advisable
to use a cloud provider because it comes with a table visualization studio.

## Testing
The test suite an be run as follows
```bash
make test
```
The following endpoints can be tested
```bash
POST : http://localhost:3000/login          # Log in and receive JWT token
GET : http://localhost:3000/account         # Fetching all acc details
POST : http://localhost:3000/account        # For creating acc 
GET : http://localhost:3000/account/{id}    # Fetching particular acc details
DELETE : http://localhost:3000/account/{id} # Deleting particular acc
POST : http://localhost:3000/transfer       # Transfering money to an account
```

The header and body requirements of the endpoints can be found from the 
`types.go` file. Will share a link to the Postman collection later.
