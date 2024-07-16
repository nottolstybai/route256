package httpserver

import (
	"net/http"
	"strconv"
)

type ValidateUser struct {
	UserID int64 `validate:"required"`
}

func (s *CartServer) DeleteAllItems(w http.ResponseWriter, r *http.Request) {
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

	if err := s.cart.DeleteItemsByUserID(r.Context(), userID); err != nil {
		http.Error(w, "failed to delete items", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
