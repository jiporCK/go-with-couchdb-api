set shell := ["bash", "-cu"]
run:
    go run ./cmd/app/main.go # Run the application
 
tidy:
    go mod tidy # Tidy the project (check go version you installed and go.mod version)