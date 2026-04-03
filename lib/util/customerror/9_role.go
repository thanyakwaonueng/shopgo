package customerror

// Role errors (09xxx)
const featRoleNum = 9

var (
	// Role validation errors (0900x)
	ROLE_NAME_DUPLICATE              = New(featRoleNum, 0, "Role name already exists")                                    // 09000
	ROLE_NAME_NOT_FOUND              = New(featRoleNum, 1, "Role name not found")                                         // 09001
	ROLE_NOT_FOUND                   = New(featRoleNum, 2, "Role not found")                                              // 09002
	ROLE_CANNOT_DELETE_ASSIGNED_ROLE = New(featRoleNum, 3, "Cannot delete role because it is assigned to existing users") // 09003
)
