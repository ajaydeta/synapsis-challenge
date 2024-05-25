package service

import (
	"context"
	"github.com/google/uuid"
	errors "github.com/rotisserie/eris"
	"synapsis-challenge/internal/core/domain"
	"synapsis-challenge/internal/core/port/outbound/registry"
	"synapsis-challenge/shared"
)

type ShoppingCartService struct {
	repositoryRegistry registry.RepositoryRegistry
}

func NewShoppingCartService(repositoryRegistry registry.RepositoryRegistry) *ShoppingCartService {
	return &ShoppingCartService{
		repositoryRegistry: repositoryRegistry,
	}
}

func (s *ShoppingCartService) Add(ctx context.Context, param *domain.ShoppingCart) (string, error) {
	var (
		err              error
		id               string
		shoppingCart     *domain.ShoppingCart
		productRepo      = s.repositoryRegistry.GetProductRepository()
		shoppingCartRepo = s.repositoryRegistry.GetShoppingCartRepository()
		customerId       = ctx.Value("customerId").(string)
	)

	_, err = productRepo.FindById(ctx, param.ProductID)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return id, err
		}
		return id, errors.Wrap(err, "Add.productRepo.FindById")
	}

	shoppingCart, err = shoppingCartRepo.FindByCustomerProductId(ctx, customerId, param.ProductID)
	if err != nil && !errors.Is(err, shared.ErrNotFound) {
		return id, errors.Wrap(err, "ShoppingCart.FindByCustomerProductId")
	}

	if shoppingCart != nil {
		param.ID = shoppingCart.ID

		err = shoppingCartRepo.Update(ctx, param)
		if err != nil {
			return id, errors.Wrap(err, "ShoppingCart.Update")
		}

		return param.ID, nil
	}

	id = uuid.NewString()
	param.ID = id
	param.CustomerID = customerId

	err = shoppingCartRepo.Add(ctx, param)
	if err != nil {
		return id, errors.Wrap(err, "Add.shoppingCartRepo.Add")
	}

	return id, nil
}