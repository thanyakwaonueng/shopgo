package serviceorders

import (
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/api/service/orders/command"
	"github.com/thanyakwaonueng/shopgo/api/service/orders/query"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func Register(
	domainDb *gorm.DB,
	logger *slog.Logger,
) {
	// Register CreateOrder Command
	serviceCreateOrder := command.NewCreateOrderHandler(logger, domainDb)
	err := mediatr.RegisterRequestHandler(serviceCreateOrder)
    if err != nil {
		panic(err)
	}
    
    // Register GetOrders Query
	serviceGetOrders := query.NewGetOrdersHandler(logger, domainDb)
	err = mediatr.RegisterRequestHandler(serviceGetOrders)
    if err != nil {
		panic(err)
	}

    // Register GetOrderByID Query
	serviceGetOrderByID := query.NewGetOrderByIDHandler(logger, domainDb)
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
