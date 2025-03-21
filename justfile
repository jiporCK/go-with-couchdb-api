set shell := ["bash", "-cu"]
run:
    export COUCHDB_HOST=202.178.125.77
    export COUCHDB_USER=admin
    export COUCHDB_PASSWORD=adminpw
    export COUCHDB_DATABASE=ishopdb
    export COUCHDB_PORT=5984
    go run ./cmd/app/main.go # Run the application
 
tidy:
    go mod tidy # Tidy the project (check go version you installed and go.mod version)