package httpserver

import (
	"net/http"
	"strconv"
)

type ValidateUserSku struct {
	UserID int64 `validate:"required"`
	SkuID  int64 `validate:"required"`
}

func (s *CartServer) DeleteItem(w http.ResponseWriter, r *http.Request) {
	var err error
	validationObj := &ValidateUserSku{}

	validationObj.UserID, err = strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "failed to parse user id", http.StatusBadRequest)
		return
	}

	validationObj.SkuID, err = strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if err != nil {
		http.Error(w, "failed to parse sku id", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(validationObj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.cart.DeleteItem(r.Context(), validationObj.UserID, validationObj.SkuID); err != nil {
		http.Error(w, "failed to delete specified item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
