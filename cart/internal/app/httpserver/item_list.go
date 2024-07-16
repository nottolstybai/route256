package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"route256.ozon.ru/project/cart/internal/entity"
	"route256.ozon.ru/project/cart/internal/service"
	"strconv"
)

type Item struct {
	SkuID int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

type ListItemsResponse struct {
	Items      []Item `json:"items"`
	TotalPrice uint32 `json:"total_price"`
}

func (s *CartServer) ListItems(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "failed to parse user id", http.StatusBadRequest)
		return
	}

	reqValidationObj := ValidateUser{
		UserID: userID,
	}
	if err := validate.Struct(reqValidationObj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	itemList, err := s.cart.GetItemsByUserID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrCartEmpty) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get items", http.StatusInternalServerError)
		return
	}

	respData, err := encodeResponse(itemList)
	if err != nil {
		http.Error(w, "failed encoding response", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(respData); err != nil {
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
	return
}

func encodeResponse(itemList *entity.ListItems) ([]byte, error) {
	resp := &ListItemsResponse{}

	resp.TotalPrice = itemList.TotalPrice

	for _, item := range itemList.Items {
		resp.Items = append(resp.Items, Item{
			SkuID: item.SkuID,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	respData, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
