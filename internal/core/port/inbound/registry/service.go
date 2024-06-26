package registry

import "synapsis-challenge/internal/core/port/inbound/service"

type ServiceRegistry interface {
	GetCustomerService() service.CustomerService
	GetProductService() service.ProductService
	GetShoppingCartService() service.ShoppingCartService
	GetTransactionService() service.TransactionService
}
