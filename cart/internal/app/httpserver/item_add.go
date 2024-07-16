package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"route256.ozon.ru/project/cart/internal/client/product_service"
	"route256.ozon.ru/project/cart/internal/service"
	"strconv"
)

type AddToCartRequest struct {
	Count  uint16 `json:"count" validate:"required"`
	SkuID  int64  `validate:"required"`
	UserID int64  `validate:"required"`
}

func (s *CartServer) AddItem(w http.ResponseWriter, r *http.Request) {
	reqObject, err := getRequestData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(reqObject); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.cart.AddItem(r.Context(), reqObject.UserID, reqObject.SkuID, reqObject.Count); err != nil {
		if errors.As(err, &product_service.ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, service.ErrItemCountTooBig) {
			http.Error(w, err.Error(), http.StatusPreconditionFailed)
			return
		}
		http.Error(w, "failed to add item", http.StatusInternalServerError)
	}

	return
}

func getRequestData(r *http.Request) (req *AddToCartRequest, err error) {
	req = &AddToCartRequest{}

	req.UserID, err = strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		return nil, err
	}

	req.SkuID, err = strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil {
		return nil, err
	}

	if err = json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	return

}
