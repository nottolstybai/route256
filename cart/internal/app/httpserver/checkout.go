package httpserver

import (
	"encoding/json"
	"net/http"
)

type CheckoutRequest struct {
	User int64 `json:"user" validate:"required"`
}

type CheckoutResponse struct {
	OrderID int64 `json:"orderID"`
}

func (s *CartServer) Checkout(w http.ResponseWriter, r *http.Request) {
	var reqBody CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderID, err := s.cart.Checkout(r.Context(), reqBody.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(orderID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}
