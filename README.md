
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
