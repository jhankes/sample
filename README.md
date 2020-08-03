# sample
a sample go web app using fiber and gorm

## Run Prereqs
- Clone repo and cd into sample dir
- Docker recent version with docker-compose
- Access to Docker Hub

## Run options
1. Run without building locally with compose
```
docker-compose up
```
2. Build locally and then run with local Docker image with compose
```
docker build -t jhankes/sample:0.1 .
docker-compose up
```
3. Run locally without compose
```
docker build -t jhankes/sample:0.1 .
./config/setupLocalDockerPostgres.sh
docker run -d \
  --link sample-postgres
  --name sample \
  -p 3000:3000 \
  -e "SAMPLE_HOST=sample-postgres" \
  jhankes/sample:0.1 
```
4.  Run with an IDE, either run the db setup first or supply envs for existing db.

## Test options with app running
1.  Open collection in postman from ./test/collections, update base env
2.  Run postman collection via Docker

