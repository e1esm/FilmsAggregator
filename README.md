# Films Aggregator Project


## Description

Rest API application. This project is an API which lets users interact with the database filled with information about films. 
Users, depending on their role, may either have full access to the API:(add and delete films) or have restricted access 
but still one that serves their needs: get all the needed information about the show with multiple filters to match their query the best possible way.
There are multiple endpoints available for the interaction with the API.


## Functionality
- ### Basic:
  - Get information about all the available films
  - Get films where actors were shot in
  - Get films that were produced by a precise person
  - Get information about precise film
- ### Admins-only:
  - Add new film to the database
  - Delete film from the database

## Technolologies
- Golang 1.20
- PostgreSQL
- Reindexer
- Docker Compose
- Zap Logger
- Open API
  - Swagger (swaggo)
- Golang-migrate
- JWT 
- Tests: unit testing and integration testing (75% coverage):
  - uber-go.Mock,
  - testify,
  - Dockertest
- YML-based flexible configuration

## Required Environment:
- ### Using docker:
  - Golang 1.20
  - Docker client
- ### Without docker:
  - Golang 1.20
  - PostgresQL
  - Reindexer

## Launch:
To run the app locally you need:
- Clone this repository
- Run any command you wish from the Makefile, for example:
```
make build_dockerize
```

## Interface
- ### Interface of the Swagger UI for this project:
!["Interface of the Swagger UI](./screenshots/swagger-interface.png)
- ### List all the films in the database:
!["Listing of all the films available in the database](./screenshots/get-all.png)
- ### List all the films an actor was shot in:
!["Listing of received films](./screenshots/filter-by-actor.png)
- ### Sign up form (right now clients can assign their role themselves):
!["Sign up form](./screenshots/signup-example.png)
- ### Sign in form
!["Sign in form](./screenshots/signin-example.png)
- ### Delete film:
!["Delete film](./screenshots/delete-film-example.png)
- ### Add film:
!["Add film](./screenshots/add-film-example.png)


## Author
Egor Mikhailov\
telegram: @el_esm\
email: elesmdev1102@gmail.com