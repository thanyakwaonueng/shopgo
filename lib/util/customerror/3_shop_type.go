package customerror

// Shop Type errors (Feature 3: Shop Type) - Code: 03xxx
const featShopTypeNum = 3

var (
	SHOP_TYPE_NAME_DUPLICATE = New(featShopTypeNum, 0, "Shop type name already exists")                                                // 03000
	SHOP_TYPE_IN_USE         = New(featShopTypeNum, 1, "Cannot change active status because shop type is in use by one or more shops") // 03001
)
