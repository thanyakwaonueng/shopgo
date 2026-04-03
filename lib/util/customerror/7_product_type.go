package customerror

const featProductTypeNum = 7

var (
	PRODUCT_TYPE_DUPLICATE_NAME         = New(featProductTypeNum, 0, "Product type name already exists")     // 07000
	PRODUCT_TYPE_HAS_ASSOCIATED_PRODUCT = New(featProductTypeNum, 1, "Product type has associated products") // 07001
)
