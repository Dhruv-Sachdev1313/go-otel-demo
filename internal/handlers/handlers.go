package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"go-otel-demo/internal/config"
	"go-otel-demo/internal/models"
	"go-otel-demo/internal/telemetry"
)

// Handler holds the dependencies for all handlers
type Handler struct {
	cfg       *config.Config
	collector *telemetry.MetricsCollector
}

// New creates a new handler instance
func New(cfg *config.Config, collector *telemetry.MetricsCollector) *Handler {
	return &Handler{
		cfg:       cfg,
		collector: collector,
	}
}

// Health handles health check requests
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer(h.cfg.ServiceName)
	_, span := tracer.Start(r.Context(), "health_check")
	defer span.End()

	span.SetAttributes(attribute.String("health.status", "ok"))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

// Error simulates random errors for testing
func (h *Handler) Error(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer(h.cfg.ServiceName)
	_, span := tracer.Start(r.Context(), "simulate_error")
	defer span.End()

	// Simulate random errors
	if rand.Float32() < 0.3 {
		span.SetAttributes(
			attribute.String("error.type", "internal_server_error"),
			attribute.Bool("error", true),
		)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error")
		return
	}
	if rand.Float32() < 0.2 {
		span.SetAttributes(
			attribute.String("error.type", "bad_request"),
			attribute.Bool("error", true),
		)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Bad Request")
		return
	}

	span.SetAttributes(attribute.String("result", "success"))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Success")
}

// AddToCart adds an item to a user's cart
func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer(h.cfg.ServiceName)
	_, span := tracer.Start(r.Context(), "add_to_cart")
	defer span.End()

	userID := r.URL.Query().Get("user_id")
	item := r.URL.Query().Get("item")

	span.SetAttributes(
		attribute.String("user.id", userID),
		attribute.String("cart.item", item),
	)

	if userID == "" || item == "" {
		span.SetAttributes(attribute.Bool("error", true))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing user_id or item parameter")
		return
	}

	// Add item to cart
	if h.collector.Carts[userID] == nil {
		h.collector.Carts[userID] = &models.ShoppingCart{
			UserID: userID,
			Items:  []string{},
		}
	}

	h.collector.Carts[userID].Items = append(h.collector.Carts[userID].Items, item)

	cartSize := len(h.collector.Carts[userID].Items)
	span.SetAttributes(
		attribute.Int("cart.size", cartSize),
		attribute.String("operation", "add"),
	)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Added %s to cart for user %s. Cart now has %d items.",
		item, userID, cartSize)
}

// RemoveFromCart removes an item from a user's cart
func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer(h.cfg.ServiceName)
	_, span := tracer.Start(r.Context(), "remove_from_cart")
	defer span.End()

	userID := r.URL.Query().Get("user_id")
	itemIndex := r.URL.Query().Get("index")

	span.SetAttributes(
		attribute.String("user.id", userID),
		attribute.String("item.index", itemIndex),
	)

	if userID == "" || itemIndex == "" {
		span.SetAttributes(attribute.Bool("error", true))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing user_id or index parameter")
		return
	}

	index, err := strconv.Atoi(itemIndex)
	if err != nil {
		span.SetAttributes(attribute.Bool("error", true))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid index parameter")
		return
	}

	cart, exists := h.collector.Carts[userID]
	if !exists || len(cart.Items) == 0 {
		span.SetAttributes(attribute.Bool("error", true))
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Cart not found or empty for user %s", userID)
		return
	}

	if index < 0 || index >= len(cart.Items) {
		span.SetAttributes(attribute.Bool("error", true))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Index out of range")
		return
	}

	// Remove item from cart
	cart.Items = append(cart.Items[:index], cart.Items[index+1:]...)

	cartSize := len(cart.Items)
	span.SetAttributes(
		attribute.Int("cart.size", cartSize),
		attribute.String("operation", "remove"),
	)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Removed item at index %d from cart for user %s. Cart now has %d items.",
		index, userID, cartSize)
}

// GetCart retrieves a user's cart contents
func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer(h.cfg.ServiceName)
	_, span := tracer.Start(r.Context(), "get_cart")
	defer span.End()

	userID := r.URL.Query().Get("user_id")

	span.SetAttributes(attribute.String("user.id", userID))

	if userID == "" {
		span.SetAttributes(attribute.Bool("error", true))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing user_id parameter")
		return
	}

	cart, exists := h.collector.Carts[userID]
	if !exists {
		span.SetAttributes(attribute.Int("cart.size", 0))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Cart for user %s: 0 items", userID)
		return
	}

	cartSize := len(cart.Items)
	span.SetAttributes(attribute.Int("cart.size", cartSize))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Cart for user %s: %d items - %v", userID, cartSize, cart.Items)
}
