package product_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go.opentelemetry.io/otel"
	"net/http"
	"route256.ozon.ru/project/cart/internal/entity"
)

const maxRetries = 3

var ErrProductNotFound = errors.New("couldn't find a product with specified skuID")

type GetProductRequest struct {
	Token string `json:"token"`
	SKU   int64  `json:"sku"`
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type ProductServiceClient struct {
	host   string
	token  string
	client *http.Client
}

func NewProductServiceClient(host, token string) *ProductServiceClient {
	client := http.DefaultClient
	client.Transport = NewRetryMiddleware(http.DefaultTransport, maxRetries)
	return &ProductServiceClient{client: client, token: token, host: host}
}

func (c *ProductServiceClient) GetProduct(ctx context.Context, skuID int64) (*entity.Product, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "ProductServiceClient.GetProduct")
	defer span.End()

	body, err := json.Marshal(GetProductRequest{
		Token: c.token,
		SKU:   skuID,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.host+"/get_product", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrProductNotFound
	}

	var respData GetProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	return &entity.Product{
		Name:  respData.Name,
		Price: respData.Price,
	}, nil
}
