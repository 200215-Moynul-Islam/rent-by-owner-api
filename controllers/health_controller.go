package controllers

import (
	"rent-by-owner-api/utils"
	"net/http"
)

// HealthController handles health-related endpoints.
type HealthController struct {
	BaseController
}

func (controller *HealthController) Health() {
	utils.SendJSONResponse(
		controller.Ctx,
		http.StatusOK,
		true,
		"Rent by Owner API is running",
		nil,
	)
}