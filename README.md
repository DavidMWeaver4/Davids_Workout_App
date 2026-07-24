## David's Workout App
![Status](https://github.com/DavidMWeaver4/Davids_Workout_App/actions/workflows/ci.yml/badge.svg)

Backend REST API consisting of 49 endpoints, PostgreSQL persistence, JWT authentication, Dockerized development, and an expanding suite of unit tests.

## Features
- 49 REST API endpoints
- JWT-based authentication
- Authorization middleware
- Unit tested handlers and authentication logic
- SQLC-generated type-safe database access
- PostgreSQL database
- GitHub Actions CI pipeline
- Docker-based local development
- Goose migrations
- Workout session management
- Exercise tracking
- Weight and cardio logging
- Custom exercise creation
- Workout statistics (volume, duration, distance)


This application allows users to:
- Create accounts securely using JWT authentication
- Create and manage workout sessions
- Track exercises, sets, reps, and weight
- Create custom exercises
- Keep personal workout data private between users

The backend is built with Go and PostgreSQL. SQLC is used to generate type-safe database queries while Goose manages database migrations, and Docker provides a reproducible local development enviroment.

Technologies:
- Go
- PostgreSQL
- SQLC
- Goose migrations
- JWT Authentication
- Docker

> This project is not finished and is in active development
## Motivations
Recently, I started focusing more on fitness and wanted a simple way to plan and track workouts. Most workout apps I tried required expensive monthly subscriptions to create your own workouts and were heavily focused on AI features that felt unnecessary for basic workout tracking.

I decided to build my own workout app both as a personal tool and as a way to improve my software engineering skills.


## Getting Started
In order to run this application you need:
- Docker (developed and tested with version 4.71.0)

## Quick Start

> Using Docker
1. Clone the repo
2. Run:
  ```bash
  make up
  ```
3. To shut down:
  ```bash
  make down
  ```

> Local Dev mode
1. Requires:
    - Go (developed with Go version 1.25.6)
    - Goose
    - Copy `.env.example` to `.env` and fill in your values
2. Start the Database:
  ```bash
  make up-db
  ```
3. Run migrations:
  ```bash
  make migrate-up
  ```
4. Start the API:
  ```bash
  make run 
  ```

## Deployment Notes
For ease of testing, this project includes default local database credentials and a fallback JWT secret so the application can be run immediately after cloning.

These defaults exist only for local development and portfolio demonstration purposes.
In a real deployment environment, all configuration values should be supplied through environment variables.

## Usage
A few simple commands to start testing. 
> Emails are unique so please change if these do not work.

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
| POST | `/login` | Login and receive a JWT |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users` | List all users |
| GET | `/users/me` | Get current user |
| PUT | `/users/me/email` | Update email |
| PUT | `/users/me/password` | Update password |
| DELETE | `/users/me` | Delete account |

### Workout Sessions
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/workout_sessions` | Create a workout session |
| GET | `/workout_sessions/me` | Get all of the current user's sessions |
| GET | `/workout_sessions/{id}` | Get session by ID |
| PUT | `/workout_sessions/{id}` | Update a workout session |
| DELETE | `/workout_sessions/{id}` | Delete a workout session |
| GET | `/workout_sessions/last/me` | Get current user's last session |
| GET | `/workout_sessions/last/count/{lastX}` | Get last X workout sessions |
| GET | `/workout_sessions/count/me` | Count the current user's sessions |
| GET | `/workout_sessions/search/desc` | Search workout sessions by description |
| GET | `/workout_sessions/search/date` | Search workout sessions by date range |

### Workout Exercises
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/workout_exercises` | Add an exercise to a workout session |
| GET | `/workout_exercises/{id}` | Get a workout exercise |
| PUT | `/workout_exercises/{id}` | Update a workout exercise |
| DELETE | `/workout_exercises/{id}` | Remove exercise from session |
| GET | `/workout_sessions/{session_id}/exercises` | Get all exercises in a workout session |
| GET | `/workout_sessions/{session_id}/exercises/count` | Count exercises in a workout session |

### Weights & Sets
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/weight_sets` | Create a weight set |
| GET | `/workout_exercises/{id}/weight_sets` | Get all weight sets for an exercise |
| PUT | `/weight_sets/{id}` | Update a weight set |
| DELETE | `/weight_sets/{id}` | Delete a weight set |
| GET | `/weight_sets/{id}/volume` | Get volume for a weight set |
| GET | `/weight_sets/{id}/duration` | Get duration for a weight set |
| GET | `/workout_exercises/{id}/volume` | Get total volume across all weight sets |
| GET | `/workout_exercises/{id}/duration` | Get total duration across all weight sets |
| GET | `/workout_sessions/{id}/weight_sets` | Get all weight sets in a workout session |
| GET | `/workout_sessions/{id}/volume` | Get total workout session volume |

### Cardio & Sets
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/cardio_sets` | Create a cardio set |
| GET | `/workout_exercises/{id}/cardio_sets` | Get all cardio sets for an exercise |
| PUT | `/cardio_sets/{id}` | Update a cardio set |
| DELETE | `/cardio_sets/{id}` | Delete a cardio set |
| GET | `/cardio_sets/{id}/distance` | Get distance for a cardio set |
| GET | `/cardio_sets/{id}/duration` | Get duration for a cardio set |
| GET | `/workout_exercises/{id}/distance` | Get total distance across all cardio sets |
| GET | `/workout_exercises/{id}/duration` | Get total duration across all cardio sets |
| GET | `/workout_sessions/{id}/cardio_sets` | Get all cardio sets in a workout session |
| GET | `/workout_sessions/{id}/distance` | Get total workout session distance |
| GET | `/workout_sessions/{id}/duration` | Get total workout session duration |

### Exercises
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/exercises` | Create a custom exercise |
| DELETE | `/exercises/{id}` | Delete a custom exercise |
| GET | `/exercises/search` | Search available exercises |
| GET | `/exercises/{id}` | Get exercise by ID |
| PUT | `/exercises/{id}` | Update a custom exercise |



## TODO
- [x] Initial unit test suite
- [ ] Increase test coverage (Currently implementing)
    - /cmd/api:  42.8% coverage
    - /internal/auth: 56.0% coverage
- [ ] CD pipeline
- [ ] Basic frontend 


## Contributing

### Clone the repo

```bash
git clone https://github.com/DavidMWeaver4/Davids_Workout_App
cd Davids_Workout_App
```

### Build the compiled binary

```bash
make build
```

### Run the test suite

```bash
make test
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
