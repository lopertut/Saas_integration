package main

import (
    "backend/handler"
    "backend/middleware"
    "backend/repository"
    "backend/service"
    "context"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/handlers" // 1. Import gorilla handlers
    "github.com/gorilla/mux"
    "github.com/jackc/pgx/v5/pgxpool"
    _ "github.com/joho/godotenv/autoload"

	"encoding/json"
	"github.com/stripe/stripe-go/v79"
    "github.com/stripe/stripe-go/v79/checkout/session"
)

var (
    err error
)

func main() {
	ctx := context.Background()

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	
	// 1. Properly get the port and ensure it has a colon prefix (fallback to ":8080" if empty)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	repo := repository.NewRepository(pool)
	s := service.NewService(repo)
	h := handler.NewHandler(s)
	authService := service.NewAuthService(repo, os.Getenv("JWT_KEY"))
	authHandler := handler.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	router := mux.NewRouter()

	router.HandleFunc("/products", h.GetProducts).Methods("GET")
	router.HandleFunc("/products/{id}", h.GetProductById).Methods("GET")
	router.HandleFunc("/products/category/{name}", h.GetProductsByCategoryId).Methods("GET")
	router.HandleFunc("/products/search/{searchQuery}", h.GetProductsByName).Methods("GET")

	router.HandleFunc("/registration", authHandler.Registration).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	router.Handle("/cart", authMiddleware.Protect(http.HandlerFunc(h.GetCartItemsByCartId))).Methods("GET")
	router.Handle("/cartItem", authMiddleware.Protect(http.HandlerFunc(h.AddCartItem))).Methods("POST")
	router.HandleFunc("/cartItem/increase/{id}", h.IncreaseCartItem).Methods("PUT")
	router.HandleFunc("/cartItem/decrease/{id}", h.DecreaseCartItem).Methods("PUT")
	router.HandleFunc("/cartItem/{id}", h.DeleteCartItem).Methods("DELETE")

	router.HandleFunc("/reviews/{product_id}", h.GetReviews).Methods("GET")
	router.Handle("/reviews", authMiddleware.Protect(http.HandlerFunc(h.AddReview))).Methods("POST")

	router.Handle("/order", authMiddleware.Protect(http.HandlerFunc(h.CreateOrder))).Methods("POST")
	router.Handle("/order", authMiddleware.Protect(http.HandlerFunc(h.GetOrders))).Methods("GET")
	router.HandleFunc("/orderItems/{id}", h.GetOrderItemsByOrderId).Methods("GET")

	router.Handle("/create-checkout-session", http.HandlerFunc(createCheckoutSession)).Methods("POST")

	// Define CORS options
	corsObj := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Accept"}),
	)

	fmt.Printf("server is running on: http://localhost%s\n", port)

	// Wrap your router with corsObj inside ListenAndServe
	log.Fatal(http.ListenAndServe(port, corsObj(router)))
}



type CheckoutRequest struct {
    Total float64 `json:"total"`
}

func createCheckoutSession(w http.ResponseWriter, r *http.Request) {
    var req CheckoutRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Total <= 0 {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Convert total to cents (Stripe expects amounts in the smallest currency unit)
    amountInCents := int64(req.Total * 100)

    params := &stripe.CheckoutSessionParams{
        PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
        LineItems: []*stripe.CheckoutSessionLineItemParams{
            {
                PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
                    Currency: stripe.String("usd"),
                    ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
                        Name: stripe.String("NEXORA Order"),
                    },
                    UnitAmount: stripe.Int64(amountInCents),
                },
                Quantity: stripe.Int64(1),
            },
        },
        Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
        SuccessURL: stripe.String("http://localhost:8000/index.html?success=true"),
        CancelURL:  stripe.String("http://localhost:8000/index.html?canceled=true"),
    }

    s, err := session.New(params)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"url": s.URL})
}
