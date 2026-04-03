package customerror

// Menu errors (Feature 10: Menu) - Code: 10xxx
const featMenuNum = 10

var (
	MENU_ROLE_ID_NOT_FOUND = New(featMenuNum, 0, "Role ID not found") // 10000
)
