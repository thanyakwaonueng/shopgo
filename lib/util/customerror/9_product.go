package customerror

const featProductNum = 9

var (
	PRODUCT_DUPLICATE_NAME = New(featProductNum, 0, "Product name already exists") // 09000
	PRODUCT_DUPLICATE_SKU  = New(featProductNum, 1, "Product SKU already exists")  // 09001
	PRODUCT_TYPE_NOT_FOUND = New(featProductNum, 2, "Product type not found")      // 09002
)
