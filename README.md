David's Workout App
>This project is not finish and is in active development

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
Golang / PostgreSQL / Sqlc
Docker / JWT / Goose

In order to run this application you need:
Golang 1.25.6
Docker (I built on version 4.71.0)


## Getting Started

1. Clone the repo
2. Copy `.env.example` to `.env` and fill in your values
3. Run `make help` to see all available commands

### Quick Start
```bash
make up          # start the database
make migrate-up  # run migrations
make run         # start the API
```
A few simple commands to test.

```
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



>TODO List
exercise handlers
weights and sets handlers
cardio handlers
go unit tests
CI/CD tests
