package customerror

const featTerritoryNum = 10

var (
	TERRITORY_DUPLICATE_NAME     = New(featTerritoryNum, 0, "Territory name already exists")                               // 10000
	TERRITORY_WAREHOUSE_IN_USE   = New(featTerritoryNum, 1, "Warehouse is already used by another territory")              // 10001
	TERRITORY_SALESPERSON_IN_USE = New(featTerritoryNum, 2, "Salesperson is already assigned to another active territory") // 10002
	TERRITORY_VEHICLE_IN_USE     = New(featTerritoryNum, 3, "Vehicle is already assigned to another active territory")     // 10003
)
