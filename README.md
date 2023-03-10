# gh-scanner

gh-scanner is a lightweight web service that provides Github code scan tool to find any potentially sensitive information that may have been accidentally committed and exposed to the public. It's designed to be simple by utilizing a scanner to search for particular keyword in the codebase.

## Overview
![gh-scanner](docs/images/gh-scanner-service.png)

## Features
- Scan for potential secrets and credentials and store findings in DB
- Use message queue as a channel to publish/subscribe messages for code scan request
- REST API for CRUD operation on the Github repository
- REST API for CRUD operation on the scan result


## Components
### 1. API Service
The API service is the web service responsible for exposing REST API to the end user, handling all application logic, configuring and injecting required libraries, and producing Github scan requests to the message queue to trigger the scanner service to start scanning. 
### 2. Scanner Service
The scanner service is mainly used for code scanning and detecting potential secrets and credentials. It works by asynchronously searching for a set of a particular words defined in the user-provided [scanner rules](/app/worker/rules.go). The scanner service subscribes to a message channel and waits for the request to scan the Github repo sent by the API service.




# Getting Started

## Building and Running

### 1. Prerequisite
If it's your first time running this application, you need to setup app configuration with:
```
make setup
```
This will setup `.env` and configure docker network used by all gh-scanner containers for you.

Then you need to have app dependencies available, you will need Kafka with Zookeeper (as a message queue), and PostgreSQL (as a database) installed and running. Alternatively, you can use the below command to build and start Docker containers for these prerequisite services:
```
make start-service
```

**IMPORTANT**: Please note that before going next step, you need to ensure that all of the above services start and run successfully as some dependencies (i.e., Kafka) may take some time to start. In case you experience an issue with installation or some of the service not running, please try to stop services with the following command and start them again:
```
make stop-service
```
In addition, you can troubleshoot the cause of the issue by checking container log with commands `docker logs $CONTAINER_NAME`


### 2. Running gh-scan
In order to start the app, there are two services that need to be run together ??? [API service](#1-api-service) and [Scanner service](#2-running-gh-scan). To setup, please run command
```
make start-app
```

Now, you should have 5 containers running (can inspect running docker containers with `docker ps`):
1. broker (kafka)
2. ghs-zookeeper
3. ghs-db (postgres)
4. gh-scanner-api (api service)
5. gh-scanner-worker (scanner service)

 You can now start sending a request to `http://localhost:8080`

This project also includes with [Postman collection](/gh-scanner.postman_collection.json). Please download and import to your machine and set Postman environment variable `gh_scanner_host` to `http://localhost:8080`.



## Unit test
The project comes with unit tests to ensure the validity of functionalities. In particular, all program use cases have accompanying unit tests for all scenarios in that use case. You can run the unit test with the command:
```
make test
```
This will trigger all unit tests in the project.


## Development
To run application locally, you need to first install project Go dependencies. At the root of project, run:
```
go mod download
```

Then you need to start start API Service and Scanner worker. To start API service please run:
```
go run cmd/api/main.go
```
Optional parameters:
- `--production`: read and use production environment variables 
-  `--check-migration`: check and run migration automatically if necessary

Then to start Scanner service, open a new terminal and run:
```
go run cmd/worker/main.go
```

### Run migration
If you run app locally for the first time, you may need to run the migration to create all tables with the following command:
```
go run cmd/api/main.go migrate
```

## Project Structures
The project is made easy for maintenance and for future changes by implementing clean architecture where the code is organized as a layer with a specific responsibility. The layer in this project can be defined from the innermost layer to the outermost layer as follows:
- Presenter (`app/presenters`): receive incoming requests and pass them to use case, format response sent by use case to end user
- Use cases (`app/usecases`): handle code business logic
- Entities (`app/entities`): define the data model used in the app and connect to the database