package handler

import (
	"net/http"

	"backend/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"strconv"
)

func (h *Handler) GetReviews(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	strId := vars["product_id"]

	productId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	reviews, err := h.service.GetReviews(ctx, int(productId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func (h *Handler) AddReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value("user_id").(int)
	var review model.Review

	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		log.Printf("%v", err)
		return
	}

	review.UserId = userId

	err = h.service.AddReview(ctx, review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("%v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("review added successfully"))
}
