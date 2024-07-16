package service

import (
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"route256.ozon.ru/project/cart/internal/entity"
	"route256.ozon.ru/project/cart/pkg/logger"
	"route256.ozon.ru/project/cart/pkg/synchronization/errgroup"
)

var (
	ErrCartEmpty       = errors.New("cart is empty")
	ErrItemCountTooBig = errors.New("error available count of stock lower than requested count")
)

//go:generate minimock [-i Repository, Client, LomsClient] [-o mock_repo_test.go, mock_client_test.go, mock_loms_client_test.go]
type Repository interface {
	Insert(ctx context.Context, userID, skuID int64, count uint16) error
	Remove(ctx context.Context, userID, skuID int64) error
	RemoveByUserID(ctx context.Context, userID int64) error
	List(ctx context.Context, userID int64) ([]entity.ProductInfo, error)
}

type Client interface {
	GetProduct(ctx context.Context, skuID int64) (*entity.Product, error)
}

type LomsClient interface {
	OrderCreate(ctx context.Context, userID int64, itemList *entity.ListItems) (int64, error)
	StocksInfo(ctx context.Context, sku uint32) (uint16, error)
}

type CartService struct {
	lomsService    LomsClient
	productService Client
	repository     Repository
}

func NewCartService(productService Client, repository Repository, lomsService LomsClient) *CartService {
	return &CartService{
		lomsService:    lomsService,
		productService: productService,
		repository:     repository,
	}
}

func (c *CartService) AddItem(ctx context.Context, userID, skuID int64, count uint16) error {
	ctx, span := otel.Tracer("default").Start(ctx, "service.AddItem")
	defer span.End()

	_, err := c.productService.GetProduct(ctx, skuID)
	if err != nil {
		return logger.WithError(err, "failed getting product")
	}

	availableCount, err := c.lomsService.StocksInfo(ctx, uint32(skuID))
	if err != nil {
		return logger.WithError(err, "failed getting available count of product")
	}

	if availableCount < count {
		return logger.WithError(ErrItemCountTooBig, ErrItemCountTooBig.Error())
	}

	if err := c.repository.Insert(ctx, userID, skuID, count); err != nil {
		return logger.WithError(err, "failed inserting item into repository")
	}

	logger.Info("Item added to cart",
		zap.Int64("userID", userID),
		zap.Int64("skuID", skuID),
		zap.Uint16("count", count))
	return nil
}

func (c *CartService) DeleteItem(ctx context.Context, userID, skuID int64) error {
	ctx, span := otel.Tracer("default").Start(ctx, "service.DeleteItem")
	defer span.End()

	if err := c.repository.Remove(ctx, userID, skuID); err != nil {
		return logger.WithError(err, "failed removing item from repository")
	}

	logger.Info("Item deleted from cart",
		zap.Int64("userID", userID),
		zap.Int64("skuID", skuID))
	return nil
}

func (c *CartService) DeleteItemsByUserID(ctx context.Context, userID int64) error {
	ctx, span := otel.Tracer("default").Start(ctx, "service.DeleteItemsByUserID")
	defer span.End()

	if err := c.repository.RemoveByUserID(ctx, userID); err != nil {
		return logger.WithError(err, "failed removing all items from repository")
	}

	logger.Info("All items deleted from cart", zap.Int64("userID", userID))
	return nil
}

func (c *CartService) Checkout(ctx context.Context, userID int64) (int64, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "service.Checkout")
	defer span.End()

	items, err := c.GetItemsByUserID(ctx, userID)
	if err != nil {
		return 0, logger.WithError(err, "failed listing items from cart")
	}

	orderID, err := c.lomsService.OrderCreate(ctx, userID, items)
	if err != nil {
		return 0, logger.WithError(err, "failed creating order")
	}

	if err := c.DeleteItemsByUserID(ctx, userID); err != nil {
		return 0, logger.WithError(err, "failed deleting all items from cart")
	}
	return orderID, nil
}

func (c *CartService) GetItemsByUserID(ctx context.Context, userID int64) (*entity.ListItems, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "service.GetItemsByUserID")
	defer span.End()

	productsInCart, err := c.repository.List(ctx, userID)
	if err != nil {
		return nil, logger.WithError(err, "failed getting items from repository")
	}

	if productsInCart == nil {
		return nil, logger.WithError(ErrCartEmpty, "cart is empty")
	}

	cart, err := c.fetchCartInfo(ctx, productsInCart)
	if err != nil {
		return nil, logger.WithError(err, "failed fetching data from cart")
	}

	logger.Info("Items retrieved from cart", zap.Int64("userID", userID))
	return cart, nil
}

func (c *CartService) fetchCartInfo(ctx context.Context, productsInCart []entity.ProductInfo) (*entity.ListItems, error) {
	var (
		cart       entity.ListItems
		totalPrice uint32
	)

	g, errCtx := errgroup.WithContext(ctx, 10)

	type productAndInfo struct {
		product *entity.Product
		count   uint16
		skuID   int64
	}
	productsCh := make(chan productAndInfo, len(productsInCart))

	for _, p := range productsInCart {
		g.Go(func() error { // concurrently add products info to chan
			product, err := c.productService.GetProduct(errCtx, p.SkuID)
			if err != nil {
				return err
			}
			productsCh <- productAndInfo{product: product, count: p.Count, skuID: p.SkuID}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(productsCh)

	for p := range productsCh { // read data from chan and generate response
		totalPrice += p.product.Price * uint32(p.count)

		cart.Items = append(cart.Items, entity.Item{
			SkuID: p.skuID,
			Count: p.count,
			Name:  p.product.Name,
			Price: p.product.Price,
		})
	}
	cart.TotalPrice = totalPrice
	return &cart, nil
}
