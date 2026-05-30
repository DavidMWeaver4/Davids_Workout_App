## David's Workout App
>This project is not finished and is in active development

## Motivations
Recently, I started focusing more on fitness and wanted a simple way to plan and track workouts. Most workout apps I tried required expensive monthly subscriptions to create your own workouts and were heavily focused on AI features that felt unnecessary for basic workout tracking.

I decided to build my own workout app both as a personal tool and as a way to improve my software engineering skills.


This application allows users to:
- Create accounts securely using JWT authentication
- Create and manage workout sessions
- Track exercises, sets, reps, and weight
- Create custom exercises
- Keep personal workout data private between users

The backend is built with Go and PostgreSQL. SQLC is used to generate type safe database queries, and Goose handles database migrations. The project is designed with Docker based deployment in mind.

Tools used in project:
- Golang / PostgreSQL / Sqlc
- Docker / JWT / Goose

In order to run this application you need:
- Golang 1.25.6
- Docker (I built on version 4.71.0)



## Getting Started

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
>Emails are unqiue so please change if these do not work.

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

### API Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | /api/v1/register | Create account | No |
| POST | /api/v1/login | Login | No |
| GET | /api/v1/users/me | Get profile | Yes |
| POST | /api/v1/workout_sessions | Create session | Yes |

## Contributing

### Clone the repo

```bash
git clone https://github.com/DavidMWeaver4/Davids_Workout_App
cd workout_app
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

## TODO
- [x] Weights and sets handlers
- [ ] Cardio handlers
- [ ] Unit tests
- [ ] CI/CD pipeline
