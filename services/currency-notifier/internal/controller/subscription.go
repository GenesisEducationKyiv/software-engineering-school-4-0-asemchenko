package controller

import (
	"currency-notifier/internal/exceptions"
	"currency-notifier/internal/models"
	"currency-notifier/internal/service"
	"currency-notifier/internal/util"
	"errors"
	"github.com/google/uuid"
	"net/http"
)

type SubscriptionController struct {
	Service *service.UserService
}

func NewSubscriptionController(service *service.UserService) *SubscriptionController {
	return &SubscriptionController{Service: service}
}

// Subscribe adds a new email to receive rate updates
// @Summary Subscribe to rate change notifications
// @Description Request adds a new email to receive USD to UAH exchange rate updates
// @Tags subscription
// @Accept x-www-form-urlencoded
// @Produce json
// @Param email formData string true "TargetEmail address"
// @Success 200 {string} string "TargetEmail added"
// @Failure 409 {string} string "Return if email already exists in the database"
// @Router /subscribe [post]
func (c *SubscriptionController) Subscribe(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "TargetEmail is required", http.StatusBadRequest)
		return
	}

	user, err := c.Service.Subscribe(email)
	renderSubscriptionResponse(w, user, err)
}

// Unsubscribe removes an email from receiving rate updates
// @Summary Unsubscribe from rate change notifications
// @Description Request removes an email from receiving USD to UAH exchange rate updates
// @Tags subscription
// @Accept x-www-form-urlencoded
// @Produce json
// @Param userId formData string true "Subscriber User ID"
// @Success 200 {string} string "Subscription removed"
// @Router /unsubscribe [post]
func (c *SubscriptionController) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userId")
	if userId == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = c.Service.Unsubscribe(parsedUserId)
	if err != nil {
		util.RespondJSON(w, http.StatusInternalServerError, map[string]string{
			"message": "Internal server error",
			"cause":   err.Error(),
		})
		return
	}

	util.RespondJSON(w, http.StatusOK, map[string]string{"message": "Successfully unsubscribed"})
}

func renderSubscriptionResponse(w http.ResponseWriter, user *models.User, subscribeError error) {
	if subscribeError != nil {
		if errors.Is(subscribeError, exceptions.ErrEmailAlreadyExists) {
			util.RespondJSON(w, http.StatusConflict, map[string]string{"message": "Email already subscribed"})
		} else {
			util.RespondJSON(w, http.StatusInternalServerError, map[string]string{
				"message": "Internal server error",
				"cause":   subscribeError.Error(),
			})
		}
		return
	}

	util.RespondJSON(w, http.StatusOK, map[string]string{
		"userId":  user.ID.String(),
		"message": "Subscription successful",
	})
}
