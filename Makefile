.PHONY : 

include .env
export

run:    
	go run main.go 

test:    
	go test ./... -race

