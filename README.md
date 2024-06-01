# Muzz | Technical Test

## Table of Contents

1. [Project Overview](#project-overview)
2. [Deliverables](#deliverables)
3. [Implementation Details and Assumptions](#implementation-details-and-assumptions)
4. [Prerequisites](#prerequisites)
5. [Installation](#installation)
6. [Exposed Endpoints](#exposed-endpoints)
7. [Project Structure](#project-structure)
8. [API Documentation](#api-documentation)

## Project Overview

This project is a simple dating API service built using Go. It allows users to create accounts, log in, and discover other users with age, and gender filters. while sorting by distance and attractivness score

## Deliverables

* User account creation
* User login
* Discover users with age and gender filters
* Discover users sorted by distance and attractiveness score

## Implementation Details and Assumptions

### Location Data

* User location data is assumed to be provided in latitude and longitude format

* The distance calculation between users is performed using a simplistic model, the Haversine formula

* For the purposes of this demo application, user location data is randomly generated, which means users may see potential matches from hundreds of kilometers away

* In a production environment, accurate user location data would need to be obtained through other means and sent by the frontend client

### Date of Birth / Age

* I was unsure weather I should collect the date of birth and return the calculated age in the response, as the requirement states that the /user/create endpoint should return the "date of birth" field, but the expected response shows the "age" field instead.

* I made the decision that it would be simpler to collect the user's age during account creation instead of calculating it from the date of birth. However, it's important to note that this approach is not a recommended and should be avoided in favor of using the date of birth for accurate age calculations.

### Data Storage

* For this project, a SQL database (MySQL) was chosen for its simplicity. While the current database design might not be optimal for scaling to millions of users, it serves the purpose of this demo application and provides the required functionality

### Sorting and Filtering

* The discover endpoint supports filtering by min age and max age range and gender, and sorting by distance and attractiveness score

### Attractiveness Score

* The attractiveness score is an integer value between 0 and 1, initialized to 0 during user creation

* The attractiveness score is calculated using the following formula: Attractiveness Score = (Total Likes / (Total Likes + Total Dislikes))

* New users will have a score at first, and they will be at a disadvantage, and that's okay for the purpose of this demo considering the simplicity of the approach

* The more likes a user receives, the higher their attractiveness score

### Testing

* Due to time limitations, no test files were included in the current implementation, although I acknolowdge the importance of having a thoroguh unit and integration tests

### Security

* It is assumed that the frontend client is validating the payload required for the endpoints to run correctly, and no additional input validation is implemented on the server-side.

* Assuming the frontend client is validating the payload required for the endpoints to run correctly

* Passwords are stored securely using hashing.

### User Authentication

* Token-based authentication was implemented for user authentication due to its simplicity and compatibility with the project's requirements

* Tokens are generated and securely stored in the database upon successful user login

* Tokens are required to be sent in the Authorization header for protected endpoints

### Configuration

* Configuration settings are managed through environment variables that are injected into the Docker container, allowing for easy modification when deploying to different environments like staging/production.

## Prerequisites

Before running this application, please ensure you have the following installed:

* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)

## Installation

1. Clone the repository:

    HTTPS
    ```bash
    git clone https://github.com/Eyadkht/muzz-dating.git
    ```

    or SSH
    
    ```bash
    git clone git@github.com:Eyadkht/muzz-dating.git
    ```

2. Navigate to the project directory:

    ```bash
    cd muzz-dating
    ```

3. Run the Docker script to build and start the containers:

    ```bash
    ./docker_run_dev.sh
    ```

   The script will build the Docker images the database and the backend, and start the containers

## Exposed Endpoints

The application exposes the following core endpoints:

### Create User

* http://localhost:8888/user/create

    This endpoint is responsible for creating a new user account

### Login

* http://localhost:8888/login

    This endpoint handles user authentication and returns a token for authenticated users

### Discover

* http://localhost:8888/discover

    This endpoint retrieves potential matches for the authenticated user

### Swipe

* http://localhost:8888/swipe

    This endpoint allows an authenticated user to swipe (YES or NO) on another user's profile and handles the matching logic

Detailed documentation for each endpoint, including request/response formats, and headers, can be found in the [API Documentation](#api-documentation) section below.

## Project Structure

This structure aims to maintain a clean separation of concerns, making it easier to understand, and maintain the codebase.

| Folder/File          | Purpose                                                                                     |
|--------------------- |---------------------------------------------------------------------------------------------|
| *cmd                  | Contains the application's entry point (`main.go`)  |
| *pkg                  | Contains the application's core logic                 |
| __ core              | Core functionalities like configuration, database handling, etc.         |
| __ handlers          | HTTP request handlers responsible for processing incoming requests. |
| __ models            | Data models representing the entities used in the application.                  |
| __ routes            | Defines the routes and associated handlers for different endpoints of the application.       |
| __ tests             | Directory for unit, integration, etc. (as mentioned no tests were implemented due to time limitations)                                   |
| __ utils             | Utility functions and helpers that can be used across the application.                        |
| Dockerfile           | Instructions for building a Docker image for the Go application.                                 |
| docker-compose.yml | Configuration file for Docker Compose, defining the api and database, networks, and volumes.              |
| docker_run_dev.sh  | Shell script for running the application in a development environment.                        |
| go.mod              | Used by Go modules to manage dependencies.                                               |
| go.sum              | The cryptographic checksums of the content of specific module versions. |

## API Documentation

## Create User

### Endpoint

POST /user/create

### Description

Creates a new user account. User's location data is generated randomly from the backend for this demo app.

### Request Body

| Field    | Type    | Description        |
|----------|---------|--------------------|
| email (required) | string  | User's email       |
| password (required) | string  | User's password    |
| name     (required) | string  | User's name        |
| gender   (required) | string  | User's gender (Male, Female)     |
| age      (required) | string     | User's age     |

### Example

```bash
curl -X POST \
  http://localhost:8888/user/create \
  -H 'Content-Type: application/json' \
  -d '{
        "email": "user@example.com",
        "password": "password123",
        "name": "John Doe",
        "gender": "male",
        "age": 25
    }'
```

### Responses

#### **201 Created** - User created successfully

```json
{
  "id": 123,
  "email": "example@example.com",
  "password": "hashed_password",
  "name": "John Doe",
  "gender": "male",
  "age": 30
}
```

#### **400 Bad Request** - Error decoding request body

```json
{
    "error": {
        "statusCode": 400,
        "message": "Error decoding request body: <error-message>"
    }
}
```

#### **409 Conflict** - User with this email address already exists

```json
{
    "error": {
        "statusCode": 409,
        "message": "User with this email address already exists: <email-address>"
    }
}
```

#### **500 Internal Server Error** - Error creating user

```json
{
    "error": {
        "statusCode": 500,
        "message": "Error creating user"
    }
}
```

## User Login

### Endpoint

POST /login

### Description

Authenticates user credentials and returns a token that should be used to authenticate protected endpoints.

### Request Body

| Field    | Type    | Description        |
|----------|---------|--------------------|
| email (required)   | string  | User's email       |
| password (required) | string  | User's password    |

### Example 

```bash
curl -X POST \
  http://localhost:8888/login \
  -H 'Content-Type: application/json' \
  -d '{
        "email": "user@example.com",
        "password": "password123"
    }'
```

### Responses

#### **200 OK** - User authenticated successfully

```json
{
    "token": "<generated-token>"
}
```

#### **400 Bad Request** - Error decoding request body

```json
{
    "error": {
        "statusCode": 400,
        "message": "Error decoding request body: <error-message>"
    }
}
```

#### **401 Unauthorized** - Invalid credentials

```json
{
    "error": {
        "statusCode": 401,
        "message": "Invalid credentials"
    }
}
```

#### **500 Internal Server Error** - Error retrieving user or generating user token

```json
{
    "error": {
        "statusCode": 500,
        "message": "Error retrieving user OR Error generating user Token"
    }
}
```

## Discover

### Endpoint

POST /discover

### Description

The endpoint retrieves potential matches for the authenticated user, allowing them to filter results based on minimum age, maximum age, and gender. Additionally, it excludes profiles of users on whom the authenticated user has already swiped in the past.

### Query Parameters

| Parameter    | Type    | Description        |
|----------|---------|--------------------|
| minAge (optional)   | int  | Minimum age for potential matches       |
| maxAge (optional) | int  |  Maximum age for potential matches    |
| gender (optional) | string  |  Gender of potential matches (male, female)    |

### Request Headers

| Parameter    | Value        |
|----------|---------|
| Authorization    | Token "user-token"  |

### Example 

```bash
curl -X GET \
  'http://localhost:8888/discover?minAge=20&maxAge=30&gender=female' \
  -H 'Authorization: Token <token>'
```

### Responses

#### **200 OK** - Successful retrieval of potential matches

```json
{
  "results": [
    {
      "id": 123,
      "name": "John Doe",
      "gender": "male",
      "age": 30,
      "distanceFromMe": 10.8,
      "attractivenessScore": 90.0
    },
    {
      "id": 456,
      "name": "Jane Smith",
      "gender": "female",
      "age": 25,
      "distanceFromMe": 8.2,
      "attractivenessScore": 45.0
    }
  ]
}
```

#### **401 Unauthorized** - Missing or invalid authentication token or header

```json
{
    "error": {
        "statusCode": 401,
        "message": "Invalid token"
    }
}
```

```json
{
    "error": {
        "statusCode": 401,
        "message": "Invalid Authorization header format"
    }
}
```

```json
{
    "error": {
        "statusCode": 401,
        "message": "Missing Authorization header"
    }
}
```

#### **500 Internal Server Error** - Error fetching potential matches

```json
{
    "error": {
        "statusCode": 500,
        "message": "Error fetching users"
    }
}
```

## Swipe

### Endpoint

POST /swipe

### Description

Allows a user to swipe YES or NO on another user's profile.

### Request Body

| Parameter    | Type    | Description        |
|----------|---------|--------------------|
| targetID    | int  | ID of the target user       |
| swipeType | string  | Type of swipe (YES or NO)   |

### Request Headers

| Parameter    | Value        |
|----------|---------|
| Authorization    | Token "user-token"  |

### Example

```bash
curl -X POST \
  http://localhost:8888/swipe \
  -H 'Authorization: Token <token>' \
  -H 'Content-Type: application/json' \
  -d '{
        "targetID": 123,
        "swipeType": "YES"
    }'
```

### Responses

#### **200 OK** - Swipe operation completed successfully and there was a match

```json
{
  "matched": true,
  "matchID": 123
}
```

#### **200 OK** - Swipe operation completed successfully but no match occured

```json
{
  "matched": false
}
```

#### **400 Bad Request** - Error decoding request body

```json
{
    "error": {
        "statusCode": 400,
        "message": "Error decoding request body: <error-message>"
    }
}
```

#### **400 Bad Request** - Can not swipe on yourself

```json
{
    "error": {
        "statusCode": 400,
        "message": "Cannot swipe on yourself"
    }
}
```

#### **400 Bad Request** - Match already exists

```json
{
    "error": {
        "statusCode": 400,
        "message": "Match already exists"
    }
}
```

#### **401 Unauthorized** - Missing or invalid authentication token or header

```json
{
    "error": {
        "statusCode": 401,
        "message": "Invalid token"
    }
}
```

```json
{
    "error": {
        "statusCode": 401,
        "message": "Invalid Authorization header format"
    }
}
```

```json
{
    "error": {
        "statusCode": 401,
        "message": "Missing Authorization header"
    }
}
```

#### **404 Not Found** - Target user not found

```json
{
    "error": {
        "statusCode": 404,
        "message": "Target user not found"
    }
}
```

#### **500 Internal Server Error** - Error creating swipe record

```json
{
    "error": {
        "statusCode": 500,
        "message": "Error creating swipe record"
    }
}
```
