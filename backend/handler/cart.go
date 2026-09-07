package handler

import (
	"backend/model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) GetCartItemsByCartId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value("user_id").(int)
	cartId := ctx.Value("cart_id").(int)

	println(userId)

	cartItems, err := h.service.GetCartItemsByCartId(ctx, cartId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartItems)
}

func (h *Handler) AddCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cartId := ctx.Value("cart_id").(int)
	var cartItem model.CartItem

	err := json.NewDecoder(r.Body).Decode(&cartItem)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		log.Printf("%v", err)
		return
	}

	cartItem.CartId = cartId

	err = h.service.AddCartItem(ctx, cartItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("%v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("item added successfully"))
}

func (h *Handler) IncreaseCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		http.Error(w, "invalid cart_item id", http.StatusBadRequest)
		return
	}

	err = h.service.IncreaseCartItem(ctx, int(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("item increased successfully"))
}

func (h *Handler) DecreaseCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		http.Error(w, "invalid cart_item id", http.StatusBadRequest)
		return
	}

	err = h.service.DecreaseCartItem(ctx, int(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("item decreased successfully"))
}

func (h *Handler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		http.Error(w, "invalid cart_item id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCartItem(ctx, int(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("item deleted successfully"))
}
