package controller

import (
	"currency-notifier/internal/exceptions"
	"currency-notifier/internal/service"
	"currency-notifier/internal/util"
	"errors"
	"net/http"
)

type SubscriptionController struct {
	Service *service.SubscriptionService
}

func NewSubscriptionController(service *service.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{Service: service}
}

// Subscribe adds a new email to receive rate updates
// @Summary Subscribe to rate change notifications
// @Description Request adds a new email to receive USD to UAH exchange rate updates
// @Tags subscription
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "Email address"
// @Success 200 {string} string "Email added"
// @Failure 409 {string} string "Return if email already exists in the database"
// @Router /subscribe [post]
func (c *SubscriptionController) Subscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	err := c.Service.Subscribe(email)
	if err != nil {
		if errors.Is(err, exceptions.ErrEmailAlreadySubscribed) {
			util.RespondJSON(w, http.StatusConflict, map[string]string{"message": "Email already subscribed"})
		} else {
			util.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	util.RespondJSON(w, http.StatusOK, map[string]string{"message": "Subscription successful"})
}
