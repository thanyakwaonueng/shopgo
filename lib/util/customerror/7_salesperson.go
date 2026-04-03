package customerror

const featSalespersonNum = 7

var (
	SALESPERSON_DUPLICATE_NAME        = New(featSalespersonNum, 0, "Salesperson name already exists")        // 07000
	SALESPERSON_DUPLICATE_EMPLOYEE_ID = New(featSalespersonNum, 1, "Salesperson employee ID already exists") // 07001
)
