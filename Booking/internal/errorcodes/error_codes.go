package errorcodes

type ErrorCode int

const (
	MissingEnvVarsErrorCode ErrorCode = iota
	UserNotFoundErrorCode
	GenerateJWTErrorCode
	GenerateRefreshTokenErrorCode
	UpdateLastLoginErrorCode
	FetchUserDetailsErrorCode
	FetchUserApplicationsErrorCode
	InvalidJSONFormatErrorCode
	ValidationErrorCode
	TokenExpiredErrorCode
	InvalidTokenErrorCode
	UnauthorizedErrorCode
	MissingAuthHeaderErrorCode
	InvalidAuthTypeErrorCode
	UnexpectedErrorCode
	InvalidConfigValueErrorCode
	InternalServerErrorCode
	DatabaseNotUpErrorCode
	BadRequestErrorCode
	UnprocessableEntityErrorCode
	RecordNotFoundErrorCode
	ForbiddenErrorCode
	MaximumLoginAttemptsErrorCode
)

var ErrorCodeStatus = map[ErrorCode]string{
	MissingEnvVarsErrorCode:        "TRAINING_ENV_001",
	UserNotFoundErrorCode:          "TRAINING_SQL_001",
	GenerateJWTErrorCode:           "TRAINING_TKN_001",
	GenerateRefreshTokenErrorCode:  "TRAINING_TKN_002",
	UpdateLastLoginErrorCode:       "TRAINING_SQL_002",
	FetchUserDetailsErrorCode:      "TRAINING_SQL_003",
	FetchUserApplicationsErrorCode: "TRAINING_SQL_004",
	InvalidJSONFormatErrorCode:     "TRAINING_ENV_001",
	ValidationErrorCode:            "TRAINING_VAL_001",
	BadRequestErrorCode:            "TRAINING_VAL_002",
	UnprocessableEntityErrorCode:   "TRAINING_VAL_003",
	RecordNotFoundErrorCode:        "TRAINING_REC_404",
	TokenExpiredErrorCode:          "TRAINING_TEXP_001",
	InvalidTokenErrorCode:          "TRAINING_IVAL_001",
	UnauthorizedErrorCode:          "TRAINING_AUTH_001",
	MissingAuthHeaderErrorCode:     "TRAINING_AUTH_002",
	ForbiddenErrorCode:             "TRAINING_AUTH_003",
	InvalidAuthTypeErrorCode:       "TRAINING_SQL_005",
	UnexpectedErrorCode:            "TRAINING_UNEXP_001",
	InvalidConfigValueErrorCode:    "TRAINING_CONFIG_001",
	InternalServerErrorCode:        "TRAINING_SQL_006",
	DatabaseNotUpErrorCode:         "TRAINING_UNEXP_002",
	MaximumLoginAttemptsErrorCode:  "TRAINING_MAX_ATTEMPTS_004",
}