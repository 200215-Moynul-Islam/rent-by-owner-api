package controllers

import (
	"strings"

	"rent-by-owner-api/repositories"
	"rent-by-owner-api/services"
)

type DestinationController struct {
	BaseController
	service services.DestinationService
}

func (controller *DestinationController) Prepare() {
	repository := repositories.NewDestinationRepository()
	controller.service = services.NewDestinationService(repository)
}

func (controller *DestinationController) Search() {
	query := strings.TrimSpace(
		controller.Ctx.Input.Query("q"),
	)

	if query == "" {
		controller.sendBadRequest(
			"Query parameter 'q' is required.",
		)
		return
	}

	destinations, err := controller.service.Search(query)
	if err != nil {
		controller.sendInternalServerError()
		return
	}

	controller.sendSuccess(
		"Destinations retrieved successfully.",
		destinations,
	)
}