package customerror

const featUserNum = 4

// User errors (04xxx)
var (
	USER_REQUIRED_FIELD_MISSING       = New(featUserNum, 1, "Required field is missing")                     // 04001
	USER_DUPLICATE_USERNAME           = New(featUserNum, 2, "Username already exists")                       // 04002
	USER_DUPLICATE_EMAIL              = New(featUserNum, 3, "Email already exists")                          // 04003
	USER_DUPLICATE_PHONE              = New(featUserNum, 4, "Phone number already exists")                   // 04004
	USER_NOT_FOUND                    = New(featUserNum, 5, "User not found")                                // 04005
	USER_INACTIVE                     = New(featUserNum, 6, "User is inactive")                              // 04006
	USER_ALREADY_ACTIVATED            = New(featUserNum, 7, "User is already activated")                     // 04007
	USER_CANNOT_EDIT_USERNAME         = New(featUserNum, 8, "Cannot edit username for activated user")       // 04008
	USER_CANNOT_EDIT_EMAIL            = New(featUserNum, 9, "Cannot edit email for activated user")          // 04009
	USER_NO_PERMISSION_TO_ASSIGN_ROLE = New(featUserNum, 10, "No permission to assign role for this client") // 04010
	USER_ROLE_CLIENT_MISMATCH         = New(featUserNum, 11, "User does not have access to this client")     // 04011
)
