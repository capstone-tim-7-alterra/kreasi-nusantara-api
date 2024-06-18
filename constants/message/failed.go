package message

const (
	// Password
	FAILED_HASHING_PASSWORD = "failed to hash password!"
	PASSWORD_MISMATCH       = "password mismatch!"

	// Request
	MISMATCH_DATA_TYPE   = "mismatch data type!"
	INVALID_REQUEST_DATA = "invalid request data!"

	// Database
	FAILED_CONNECT_DB = "failed connect to database: %v"
	FAILED_MIGRATE_DB = "failed to migrate database!"

	// Token
	MISSING_TOKEN           = "missing token"
	INVALID_TOKEN           = "invalid token"
	INVALID_UUID            = "invalid uuid"
	INVALID_AUTH_TYPE       = "invalid authentication type. use Bearer"
	UNAUTHORIZED            = "unauthorized user"
	FAILED_GENERATE_TOKEN   = "failed to generate token!"
	FAILED_INVALIDATE_TOKEN = "failed to invalidate token!"

	// Forbidden
	FORBIDDEN_RESOURCE = "Forbidden	resource!"

	// External Service
	EXTERNAL_SERVICE_ERROR = "External service error!"

	// User
	FAILED_CREATE_USER           = "failed to create user!"
	USER_EXIST                   = "email already exists!"
	FAILED_LOGIN                 = "login failed!"
	UNREGISTERED_EMAIL           = "unregistered email!"
	UNREGISTERED_USER            = "unregistered user!"
	DUPLICATE_KEY                = "duplicate key value violates unique constraint"
	FAILED_GET_USER              = "failed to get user!"
	FAILED_VERIFY_OTP            = "failed to verify otp!"
	USER_NOT_FOUND               = "user not found!"
	INVALID_OTP                  = "invalid otp!"
	FAILED_FORGOT_PASSWORD       = "failed to initiate forgot password!"
	FAILED_GET_PROFILE           = "failed to get user profile!"
	FAILED_UPDATE_PROFILE        = "failed to update user profile!"
	FAILED_RESET_PASSWORD        = "failed to reset password!"
	FAILED_DELETE_PROFILE        = "failed to delete user profile!"
	FAILED_GET_USER_ADDRESSES    = "failed to get user addresses!"
	FAILED_CREATE_USER_ADDRESSES = "failed to create user addresses!"
	FAILED_UPDATE_USER_ADDRESSES = "failed to update user addresses!"
	FAILED_DELETE_USER_ADDRESSES = "failed to delete user addresses!"
	FAILED_CHANGE_PASSWORD       = "failed to change password!"

	//Admin
	FAILED_CREATE_ADMIN        = "failed to create admin!"
	FAILED_LOGIN_ADMIN         = "login failed!"
	FAILED_FETCH_DATA          = "Failed to fetch data"
	FAILED_UPDATE_ADMIN        = "failed to update admin!"
	FAILED_DELETE_ADMIN        = "failed to delete admin!"
	FAILED_PARSE_ADMIN         = "failed to parse admin id"
	MISSING_USERNAME_PARAMETER = "Missing username parameter"
	FAILED_SEARCH_ADMIN        = "failed to search admin by username"
	ADMIN_NOT_FOUND            = "admin not found"
	FAILED_GET_ADMIN           = "failed to get admin!"

	//Product Admin
	FAILED_CREATE_CATEGORY        = "failed to create category"
	FAILED_PARSE_CATEGORY         = "failed to parse category id"
	FAILED_UPDATE_CATEGORY        = "failed to update category"
	FAILED_DELETE_CATEGORY        = "failed to delete category"
	FAILED_SEARCH_CATEGORY        = "failed to search category by name"
	CATEGORY_NOT_FOUND            = "category not found"
	FAILED_CREATE_PRODUCT         = "failed to create product"
	FAILED_PARSE_PRODUCT          = "failed to parse product id"
	FAILED_UPDATE_PRODUCT         = "failed to update product"
	FAILED_DELETE_PRODUCT         = "failed to delete product"
	FAILED_SEARCH_PRODUCT         = "failed to search product by name"
	PRODUCT_NOT_FOUND             = "product not found"
	INVALID_MIN_ORDER             = "invalid min order"
	INVALID_CATEGORY_ID           = "invalid category id"
	FAILED_READ_IMAGE_FILE        = "failed to read image file"
	FAILED_READ_VIDEO_FILE        = "failed to read video file"
	INVALID_PRODUCT_VARIANT_PRICE = "invalid product variant price"
	INVALID_PRODUCT_VARIANT_STOCK = "invalid product variant stock"
	INVALID_ORIGINAL_PRICE        = "invalid original price"
	INVALID_DISCOUNT_PERCENT      = "invalid discount percent"

	// Pages
	PAGE_NOT_FOUND = "page not found!"

	// Images
	FAILED_UPLOAD_IMAGE = "failed to upload image!"
	FAILED_DELETE_IMAGE = "failed to delete image!"

	// Products
	FAILED_GET_PRODUCTS        = "failed to get products!"
	FAILED_GET_PRODUCT_REVIEWS = "failed to get product reviews!"
	FAILED_CREATE_REVIEW       = "failed to create review!"
	FAILED_GET_PRODUCT_RECOMMENDATION = "failed to get product recommendation!"

	// Articles
	FAILED_GET_ARTICLES   = "failed to get articles!"
	FAILED_GET_COMMENTS   = "failed to get comments!"
	FAILED_ADD_COMMENT    = "failed to add comments!"
	FAILED_REPLY_COMMENT  = "failed to reply comment!"
	FAILED_LIKE_ARTICLE   = "failed to like article!"
	FAILED_UNLIKE_ARTICLE = "failed to unlike article!"
	FAILED_GET_ARTICLE    = "failed to get article!"

	// Events
	FAILED_GET_EVENTS = "failed to get events!"

	// Chatbot
	FAILED_ANSWER_CHAT = "maaf saya belum bisa menjawab pertanyaan yang anda ajukan. Silahkan coba pertanyaan lain"

	// Events Admin
	FAILED_CREATE_EVENTS = "failed to create events!"
	FAILED_UPDATE_EVENTS = "failed to update events!"
	FAILED_DELETE_EVENTS = "failed to delete events!"
	FAILED_SEARCH_EVENTS = "failed to search events by name"
	EVENT_NOT_FOUND      = "event not found"

	// Articles
	FAILED_CREATE_ARTICLE = "failed to create article!"
	FAILED_UPDATE_ARTICLE = "failed to update article!"
	FAILED_DELETE_ARTICLE = "failed to delete article!"
	FAILED_SEARCH_ARTICLE = "failed to search article by name"
	ARTICLE_NOT_FOUND     = "article not found!"
	NOT_FOUND             = "not found"

	// Categories
	FAILED_GET_CATEGORIES = "failed to get categories!"
	FAILED_CREATE_TICKET_TYPE = "failed to create ticket type!"
	FAILED_UPDATE_TICKET_TYPE = "failed to update ticket type!"
	FAILED_DELETE_TICKET_TYPE = "failed to delete ticket type!"
	FAILED_SEARCH_TICKET_TYPE = "failed to search ticket type by name"
	TICKET_TYPE_NOT_FOUND     = "ticket type not found"
	FAILED_GET_TICKET_TYPE    = "failed to get ticket type!"

	// Cart
	FAILED_ADD_TO_CART = "failed to add to cart!"
	FAILED_GET_CART_ITEMS = "failed to get cart items!"
	FAILED_UPDATE_CART_ITEMS = "failed to update cart items!"
	FAILED_DELETE_CART_ITEMS = "failed to delete cart items!"
)
