package customerror

// Vehicle Type errors (06xxx)
const featVehicleTypeNum = 6

var (
	// Vehicle Type creation errors (0600x)
	VEHICLE_TYPE_DUPLICATE_NAME = New(featVehicleTypeNum, 0, "Vehicle type name already exists") // 06000
)
