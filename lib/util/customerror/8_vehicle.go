package customerror

const featVehicleNum = 8

var (
	VEHICLE_DUPLICATE_NAME          = New(featVehicleNum, 0, "Vehicle name already exists")          // 08000
	VEHICLE_DUPLICATE_LICENSE_PLATE = New(featVehicleNum, 1, "Vehicle license plate already exists") // 08001
)
