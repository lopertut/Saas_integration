package handler

import (
	"encoding/json"
	"net/http"

	"backend/model"
	"github.com/gorilla/mux"
	"log"
	"strconv"
)

func (h *Handler) GetOrderItemsByOrderId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	strId := vars["id"]
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	orderItems, err := h.service.GetOrderItemsByOrderId(ctx, int(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderItems)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value("user_id").(int)

	orders, err := h.service.GetOrders(ctx, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := ctx.Value("user_id").(int)
	cartId := ctx.Value("cart_id").(int)
	var order model.Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		log.Printf("%v", err)
		return
	}

	println(userId)
	order.UserId = userId
	println(order.UserId)

	err = h.service.CreateOrder(ctx, order, cartId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("order added successfully"))
}
