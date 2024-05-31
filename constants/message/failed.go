package message

const (
	// Password
	FAILED_HASHING_PASSWORD = "failed to hash password!"
	PASSWORD_MISMATCH = "password mismatch!"

	// Request
	MISMATCH_DATA_TYPE = "mismatch data type!"
	INVALID_REQUEST_DATA = "invalid request data!"

	// Database
	FAILED_CONNECT_DB =	"failed connect to database: %v"
	FAILED_MIGRATE_DB = "failed to migrate database!"

	// Token
	MISSING_TOKEN =	"missing token"
	INVALID_TOKEN =	"invalid token"
	INVALID_AUTH_TYPE = "invalid authentication type. use Bearer"
	UNAUTHORIZED =	"unauthorized user"
	FAILED_GENERATE_TOKEN = "failed to generate token!"
	FAILED_INVALIDATE_TOKEN = "failed to invalidate token!"

	// Forbidden
	FORBIDDEN_RESOURCE = "Forbidden	resource!"

	// External Service
	EXTERNAL_SERVICE_ERROR = "External service error!"

	// User
	FAILED_CREATE_USER = "failed to create user!"
	USER_EXIST = "email already exists!"
	FAILED_LOGIN = "login failed!"
	UNREGISTERED_EMAIL = "unregistered email!"
	UNREGISTERED_USER = "unregistered user!"
	DUPLICATE_KEY = "duplicate key value violates unique constraint"
	FAILED_GET_USER = "failed to get user!"

	//Admin
	FAILED_CREATE_ADMIN = "failed to create admin!"
	FAILED_LOGIN_ADMIN = "login failed!"
	FAILED_FETCH_DATA     = "Failed to fetch data"
	FAILED_UPDATE_ADMIN = "failed to update admin!"
	FAILED_DELETE_ADMIN = "failed to delete admin!"
	FAILED_PARSE_ADMIN = "failed to parse admin id"
	MISSING_USERNAME_PARAMETER = "Missing username parameter"
	FAILED_SEARCH_ADMIN = "failed to search admin by username"
	ADMIN_NOT_FOUND = "admin not found"
	
)