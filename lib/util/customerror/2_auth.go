package customerror

// Authentication errors (Feature 2: Auth) - Code: 02xxx
const featAuthNum = 2

var (
	AUTH_INVALID_PARAMS      = New(featAuthNum, 1, "Invalid parameters")           // 02001
	AUTH_INVALID_CREDENTIALS = New(featAuthNum, 2, "Invalid username or password") // 02002

	// Token verification errors (0208x)
	AUTH_TOKEN_INVALID       = New(featAuthNum, 80, "Invalid or expired token")                    // 02080
	AUTH_TOKEN_NOT_FOUND     = New(featAuthNum, 81, "Token not found or has been used or expired") // 02081
	AUTH_TOKEN_USER_MISMATCH = New(featAuthNum, 82, "Token user mismatch")                         // 02082
	AUTH_TOKEN_ALREADY_USED  = New(featAuthNum, 83, "Token mismatch or already used")              // 02083
)
