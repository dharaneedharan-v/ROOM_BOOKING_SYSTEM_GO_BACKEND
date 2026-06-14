# Booking Service

A Go-based Room Booking Service with PostgreSQL , structured logging.

---

# Prerequisites

Before running the project, ensure the following are installed:

* Go 1.21+ (or the version used by the project)
* PostgreSQL
* Git
* Docker & Docker Compose (for containerized execution)

---

# Clone the Repository

```bash
git clone https://github.com/dharaneedharan-v/ROOM_BOOKING_SYSTEM_GO_BACKEND
cd Booking
```

---

# Install Dependencies

Install, update, and clean Go modules:

```bash
go mod tidy
```

This command:

* Downloads missing dependencies
* Removes unused dependencies
* Updates `go.sum`

---

## Configure Environment

Create/update the appropriate `.env` file with your local configuration:

The application loads environment files using:

```go
envFile := fmt.Sprintf("../../envs/.env.%s", env)
```

Examples:

```text
envs/.env.dev
envs/.env.prod
envs/.env.docker
```

```env
PORT=<your_port>

DATABASE_URL=postgresql://<your_db_user>:<your_db_password>@localhost:<your_db_port>/<your_database_name>

LOG_LEVEL=debug
LOG_FILE_MAX_SIZE_MB=10
LOG_FILE_NAME=<your_log_file_name>
LOG_DIR=./logs

PROJECT_ROOT=<your_project_root_path>

BASE_URL=http://localhost:<your_port>
```

# Running the Application Locally

Navigate to the server directory:

```powershell
cd cmd/server
```

Run the application:

```powershell
go run main.go
```

The service will start on:

```text
http://localhost:8080
```

### *FILE STRUCT*

```
├─── Booking
│   ├─── cmd
│   │   └─── server
│   ├─── deployments
│   │   ├─── docker-compose.yml
│   │   ├─── Dockerfile
│   │   └─── promtail-config.yaml
│   ├─── docs
│   │   ├─── swaggerui
│   │   │   ├─── scripts
│   │   │   │   └─── rapidoc-min.js
│   │   │   └─── index.html
│   │   └─── openapi.yaml
│   ├─── envs
│   │   ├─── .env.dev
│   │   └─── .env.docker
│   ├─── internal
│   │   ├─── config
│   │   │   └─── config.go
│   │   ├─── dtos
│   │   │   ├─── bookings_dto.go
│   │   │   └─── response_dto.go
│   │   ├─── errorcodes
│   │   │   └─── error_codes.go
│   │   ├─── handlers
│   │   │   ├─── bookings_handler.go
│   │   │   └─── routes.go
│   │   ├─── loggers
│   │   │   └─── logger.go
│   │   ├─── models
│   │   │   └─── bookings.go
│   │   ├─── repository
│   │   │   └─── bookings_repository.go
│   │   ├─── services
│   │   │   └─── bookings_service.go
│   │   └─── utils
│   │       └─── utils.go
│   ├─── pkg
│   │   ├─── database
│   │   │   ├─── database.go
│   │   │   ├─── migration.go
│   │   │   └─── seeder.go
│   │   └─── server
│   │       └─── server.go
│   ├─── tests
│   │   ├─── integration
│   │   │   └─── .gitkeep
│   │   └─── unit
│   │       ├─── config
│   │       │   ├─── mocks
│   │       │   │   └─── config_mock.go
│   │       │   └─── config_test.go
│   │       ├─── database
│   │       │   └─── database_test.go
│   │       ├─── handlers
│   │       │   ├─── routes_test.go
│   │       │   └─── sample_handler_test.go
│   │       ├─── pkg
│   │       │   └─── server
│   │       │       └─── server_test.go
│   │       ├─── repository
│   │       │   ├─── mock
│   │       │   │   └─── sample_repository_mock.go
│   │       │   └─── sample_repository_test.go
│   │       ├─── server
│   │       │   └─── server_test.go
│   │       ├─── services
│   │       │   ├─── mock
│   │       │   │   └─── sample_service_mock.go
│   │       │   └─── sample_service_test.go
│   │       └─── utils
│   │           └─── utils_test.go
│   ├─── go.mod
│   ├─── go.sum
│   └─── README.md
└─── .gitignore
```

### SAMPLES : 

<img width="1920" height="1080" alt="image" src="https://github.com/user-attachments/assets/1dd6831a-610a-44b3-aae0-e2d1fa3a2f2a" />

