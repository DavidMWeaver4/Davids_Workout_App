## David's Workout App

A RESTful workout tracking API built with Go and PostgreSQL.

## Features
- JWT-based authentication
- Workout session management
- Exercise tracking
- Weight and cardio logging
- Custom exercise creation
- Workout statistics (volume, duration, distance)
- PostgreSQL persistence
- Docker-based local development

## Motivations
Recently, I started focusing more on fitness and wanted a simple way to plan and track workouts. Most workout apps I tried required expensive monthly subscriptions to create your own workouts and were heavily focused on AI features that felt unnecessary for basic workout tracking.

I decided to build my own workout app both as a personal tool and as a way to improve my software engineering skills.


This application allows users to:
- Create accounts securely using JWT authentication
- Create and manage workout sessions
- Track exercises, sets, reps, and weight
- Create custom exercises
- Keep personal workout data private between users

The backend is built with Go and PostgreSQL. SQLC is used to generate type safe database queries, and Goose handles database migrations. The project is designed with Docker-based deployment in mind.

Tools used in project:
- Golang
- PostgreSQL
- SQLC
- Goose migrations
- JWT Authentication
- Docker

>This project is not finished and is in active development

## Getting Started
In order to run this application you need:
- Golang 1.25.6
- Docker (I built on version 4.71.0)


1. Clone the repo
2. Copy `.env.example` to `.env` and fill in your values
3. Run `make help` to see all available commands

## Quick Start
```bash
make up          # start the database
make migrate-up  # run migrations
make run         # start the API
```

## Usage
A few simple commands to test. 
>Emails are unique so please change if these do not work.

```bash
# register
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"email": "test@test.com", "password": "password123"}'

# login
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email": "test@test.com", "password": "password123"}'

# authenticated request
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer <token_from_login>"
```

## API Endpoints

Base URL: `/api/v1`

### Auth
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Register a new user |
| POST | `/login` | Login and receive tokens |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users/me` | Get current user |
| PUT | `/users/me/email` | Update email |
| PUT | `/users/me/password` | Update password |
| DELETE | `/users/me` | Delete account |

### Workout Sessions
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/workout_sessions` | Create a session |
| GET | `/workout_sessions/me` | Get all my sessions |
| GET | `/workout_sessions/{id}` | Get session by ID |
| PUT | `/workout_sessions/{id}` | Update session |
| DELETE | `/workout_sessions/{id}` | Delete session |
| GET | `/workout_sessions/last/me` | Get my last session |
| GET | `/workout_sessions/last/count/{lastX}` | Get last X sessions |
| GET | `/workout_sessions/count/me` | Count my sessions |
| GET | `/workout_sessions/search/desc` | Search by description |
| GET | `/workout_sessions/search/date` | Search by date range |

### Workout Exercises
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/workout_exercises` | Add exercise to session |
| GET | `/workout_exercises/{id}` | Get workout exercise by ID |
| DELETE | `/workout_exercises/{id}` | Remove exercise from session |
| GET | `/workout_sessions/{session_id}/exercises` | Get all exercises in session |
| GET | `/workout_sessions/{session_id}/exercises/count` | Count exercises in session |

### Weights & Sets
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/weights_and_sets` | Create a weight set |
| GET | `/workout_exercises/{id}/sets` | Get all sets for exercise |
| PUT | `/weights_and_sets/{id}` | Update a set |
| DELETE | `/weights_and_sets/{id}` | Delete a set |
| GET | `/weights_and_sets/{id}/volume` | Get volume for a set |
| GET | `/weights_and_sets/{id}/volume/all` | Get total volume across all sets |
| GET | `/weights_and_sets/{id}/duration` | Get duration for a set |
| GET | `/weights_and_sets/{id}/duration/all` | Get total duration across all sets|

### Cardio & Sets
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/cardio_and_sets` | Create a cardio set |
| GET | `/workout_exercises/{id}/cardio` | Get all cardio sets for exercise |
| PUT | `/cardio_and_sets/{id}` | Update a cardio set |
| DELETE | `/cardio_and_sets/{id}` | Delete a cardio set |
| GET | `/cardio_and_sets/{id}/distance` | Get distance for a set |
| GET | `/cardio_and_sets/{id}/distance/all` | Get total distance across all sets|
| GET | `/cardio_and_sets/{id}/duration` | Get duration for a set |
| GET | `/cardio_and_sets/{id}/duration/all` | Get total duration across all sets |

### Exercises
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/exercises` | Create a custom exercise |
| DELETE | `/exercises/{id}` | Delete a custom exercise |
| GET | `/exercises/search` | List/search available exercises |
| GET | `/exercises/{id}` | Get exercise by ID |
| PUT | `/exercises/{id}` | Update a custom exercise |



## TODO
- [x] Weights and sets handlers
- [X] Cardio handlers
- [X] Exercise handlers
- [ ] Unit tests
- [ ] CI/CD pipeline
- [ ] Frontend Client


## Contributing

### Clone the repo

```bash
git clone https://github.com/DavidMWeaver4/Davids_Workout_App
cd Davids_Workout_App
```

### Build the compiled binary

```bash
go build
```

> Sorry, test suite is still TODO

### Run the test suite

```bash
go test ./...
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.

## Bug finding

If you find a bug, please email me at [dmweaver4@gmail.com](mailto:dmweaver4@gmail.com)
