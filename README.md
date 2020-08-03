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
$ docker build -t jhankes/sample:0.1 .

...

$ docker-compose up
```
3. Run locally without compose
```
$ docker build -t jhankes/sample:0.1 .

...

$ ./config/setupLocalDockerPostgres.sh

...

$ docker run -d \
    --link sample-postgres
    --name sample \
    -p 3000:3000 \
    -e "SAMPLE_HOST=sample-postgres" \
    jhankes/sample:0.1 
```
4.  Run with an IDE, either run the db setup first or supply envs for existing db.

## Test options with app running
1.  Open collection in postman from ./test/collections, update base env
2.  Run postman collection via Docker.  This requires a little configuration based on the run option selected.  The network and container name could vary, also depends on where the repo resides to reference the newman config.  Ensure newman environment file referenced has the correct base env. Here is an example when using option 1:
```
$ docker run \
    --link sample_web_1 \
    --network sample_default \
    -v $GOPATH/src/github.com/jhankes/sample/test/data:/etc/newman \
    -t postman/newman:alpine \
    run "https://www.getpostman.com/collections/b8b5643a23d9267f5ef6" \
    -e /etc/newman/sample.postman_environment.json
```

Example results are shown below and can be [viewed here](https://github.com/jhankes/sample/blob/master/test/data/sample-tests.postman_test_run.json):

![Example test results from Newman](https://github.com/jhankes/sample/raw/master/test/data/newmandockertestresult.png)
