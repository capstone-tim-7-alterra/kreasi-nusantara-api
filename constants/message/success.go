package message

const (
	// User
	USER_CREATED_SUCCESS          = "user created successfully!"
	LOGIN_SUCCESS                 = "login success!"
	GET_USER_SUCCESS              = "user retrieve successfully!"
	OTP_SENT_SUCCESS              = "OTP sent to email!"
	PASSWORD_RESET_SUCCESS        = "password reset successfully!"
	GET_PROFILE_SUCCESS           = "user profile retrieved successfully!"
	UPDATE_PROFILE_SUCCESS        = "user profile updated successfully!"
	DELETE_PROFILE_SUCCESS        = "user profile deleted successfully!"
	GET_USER_ADRESSES_SUCCESS     = "user addresses retrieved successfully!"
	CREATE_USER_ADDRESSES_SUCCESS = "user addresses created successfully!"
	UPDATE_USER_ADDRESSES_SUCCESS = "user addresses updated successfully!"
	DELETE_USER_ADDRESSES_SUCCESS = "user addresses deleted successfully!"
	CHANGE_PASSWORD_SUCCESS       = "password changed successfully!"

	//Admin
	ADMIN_CREATED_SUCCESS   = "admin created successfully!"
	ADMIN_RETRIEVED_SUCCESS = "admin retrieve successfully!"
	ADMIN_UPDATED_SUCCESS   = "admin updated successfully!"
	ADMIN_DELETED_SUCCESS   = "admin deleted successfully!"
	SUCCESS_FETCH_DATA      = "Successfully fetched data"
	SUCCES_SEARCH_ADMIN     = "Successfully search admin"
	GET_ADMIN_SUCCESS       = "admin retrieved successfully!"

	//Products Admin
	PRODUCT_CREATED_SUCCESS  = "product created successfully!"
	CATEGORY_CREATED_SUCCESS = "category created successfully!"
	CATEGORY_UPDATED_SUCCESS = "category updated successfully!"
	CATEGORY_DELETED_SUCCESS = "category deleted successfully!"
	PRODUCT_UPDATED_SUCCESS  = "product updated successfully!"
	PRODUCT_DELETED_SUCCESS  = "product deleted successfully!"

	// Images
	UPLOAD_IMAGE_SUCCESS = "image uploaded successfully!"
	DELETE_IMAGE_SUCCESS = "image deleted successfully!"

	// Products
	GET_PRODUCTS_SUCCESS               = "products retrieved successfully!"
	GET_PRODUCT_REVIEWS_SUCCESS        = "product reviews retrieved successfully!"
	CREATE_REVIEW_SUCCESS              = "review created successfully!"
	GET_PRODUCT_RECOMMENDATION_SUCCESS = "product recommendation retrieved successfully!"

	// Articles
	GET_ARTICLES_SUCCESS   = "articles retrieved successfully!"
	GET_COMMENTS_SUCCESS   = "comments retrieved successfully!"
	ADD_COMMENT_SUCCESS    = "comment added successfully!"
	REPLY_COMMENT_SUCCESS  = "comment replied successfully!"
	LIKE_ARTICLE_SUCCESS   = "article liked successfully!"
	UNLIKE_ARTICLE_SUCCESS = "article unliked successfully!"
	GET_ARTICLE_SUCCESS    = "article retrieved successfully!"

	// Events
	GET_EVENTS_SUCCESS = "events retrieved successfully!"

	// events Admin
	CREATE_EVENTS_SUCCESS = "event created successfully!"
	UPDATE_EVENTS_SUCCESS = "event updated successfully!"
	DELETE_EVENTS_SUCCESS = "event deleted successfully!"

	// Categories Event
	GET_CATEGORY_SUCCESS    = "categories retrieved successfully!"
	UPDATE_CATEGORY_SUCCESS = "category updated successfully!"
	DELETE_CATEGORY_SUCCESS = "category deleted successfully!"
	CREATE_CATEGORY_SUCCESS = "category created successfully!"

	CREATE_TICKET_TYPE_SUCCESS = "ticket type created successfully!"
	UPDATE_TICKET_TYPE_SUCCESS = "ticket type updated successfully!"
	DELETE_TICKET_TYPE_SUCCESS = "ticket type deleted successfully!"
	GET_TICKET_TYPE_SUCCESS    = "ticket type retrieved successfully!"

	GET_PRICES        = "prices retrieved successfully!"
	GET_DETAIL_PRICES = "detail prices retrieved successfully!"
	DELETE_PRICES     = "prices deleted successfully!"
	UPDATE_PRICES     = "prices updated successfully!"

	//Articles Admin
	CREATE_ARTICLE_SUCCESS = "article created successfully!"
	UPDATE_ARTICLE_SUCCESS = "article updated successfully!"
	DELETE_ARTICLE_SUCCESS = "article deleted successfully!"

	// Cart
	ADD_TO_CART_SUCCESS       = "items added to cart successfully!"
	GET_CART_ITEMS_SUCCESS    = "items retrieved successfully!"
	UPDATE_CART_ITEMS_SUCCESS = "items updated successfully!"
	DELETE_CART_ITEMS_SUCCESS = "items deleted successfully!"
)
