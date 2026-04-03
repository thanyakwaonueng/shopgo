package customerror

const featShopNum = 5

var (
	SHOP_DUPLICATE_NAME = New(featShopNum, 0, "Shop name already exists") // 05000
)
