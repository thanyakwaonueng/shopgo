package customerror

// Middleware errors (Feature 1: Middleware) - Code: 01xxx
const featMiddlewareNum = 1

var (
	MIDDLEWARE_INVALID_TOKEN   = New(featMiddlewareNum, 2, "Invalid or missing authentication token")
	MIDDLEWARE_TOKEN_EXPIRED   = New(featMiddlewareNum, 3, "Token has expired")
	MIDDLEWARE_TOKEN_MISMATCH  = New(featMiddlewareNum, 4, "Token mismatch - session expired")
	MIDDLEWARE_TOKEN_NOT_FOUND = New(featMiddlewareNum, 5, "Token not found")
)
