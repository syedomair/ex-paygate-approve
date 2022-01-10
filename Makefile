.PHONY : 

include .env
export

run:    
	go run main.go 

rund:
	docker build -t syedomair/ex-paygate-approve:latest \
	 --build-arg service_name=$(SERVICE_NAME) \
         --build-arg log_level=$(LOG_LEVEL) \
         --build-arg port=$(PORT) \
         --build-arg database_url=$(DATABASE_URL_DOCKER) \
         --build-arg signingkey=$(SIGNINGKEY)  .
	docker container run  -p 8321:8321 --name ex-paygate-approve syedomair/ex-paygate-approve:latest

test_v:    
	go test ./... -v

test_r:    
	go test ./... -race

