package customerror

const featWarehouseNum = 6

var (
	WAREHOUSE_DUPLICATE_NAME                        = New(featWarehouseNum, 0, "Warehouse name already exists")         // 06000
	WAREHOUSE_GEOM_POINT_IS_NOT_RELATED_TO_LOCATION = New(featWarehouseNum, 1, "geom point is not related to location") // 06001
	WAREHOUSE_GEOM_POINT_IS_WITHIN_SOME_TERRITORY   = New(featWarehouseNum, 2, "geom point is within some territory")   // 06002
	WAREHOUSE_DUPLICATE_GEOM_POINT                  = New(featWarehouseNum, 3, "Warehouse geom point already exists")   // 06003
)
