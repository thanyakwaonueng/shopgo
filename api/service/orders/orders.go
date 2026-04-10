package serviceorders

import (
	"log/slog"
    repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/api/service/orders/command"
	"github.com/thanyakwaonueng/shopgo/api/service/orders/query"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func Register(
	domainDb *gorm.DB,
	logger *slog.Logger,
    repoProduct repogeneric.Product,
    repoOrder repogeneric.Order,
) {
	// Register CreateOrder Command
	serviceCreateOrder := command.NewCreateOrderHandler(logger, domainDb, repoProduct, repoOrder)
	err := mediatr.RegisterRequestHandler(serviceCreateOrder)
    if err != nil {
		panic(err)
	}
    
    // Register GetOrders Query
	serviceGetOrders := query.NewGetOrdersHandler(logger, domainDb, repoOrder)
	err = mediatr.RegisterRequestHandler(serviceGetOrders)
    if err != nil {
		panic(err)
	}

    // Register GetOrderByID Query
	serviceGetOrderByID := query.NewGetOrderByIDHandler(logger, domainDb, repoOrder)
	err = mediatr.RegisterRequestHandler(serviceGetOrderByID)
    if err != nil {
		panic(err)
	}

    // Register UpdateOrderStatus Command
	serviceUpdateStatus := command.NewUpdateOrderStatusHandler(logger, domainDb)
	err = mediatr.RegisterRequestHandler(serviceUpdateStatus)
    if err != nil {
		panic(err)
	}

    // Add CancelOrder registration
	serviceCancelOrder := command.NewCancelOrderHandler(logger, domainDb)
	err = mediatr.RegisterRequestHandler(serviceCancelOrder)
    if err != nil {
		panic(err)
	}
}
