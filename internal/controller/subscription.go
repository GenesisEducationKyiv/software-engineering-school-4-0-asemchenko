package controller

import "net/http"

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
func Subscribe(w http.ResponseWriter, r *http.Request) {

}
