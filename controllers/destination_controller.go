package controllers

import (
	"strconv"
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

func (controller *DestinationController) Autocomplete() {
	query := strings.TrimSpace(
		controller.Ctx.Input.Query("q"),
	)

	if query == "" {
		controller.sendBadRequest(
			"Query parameter 'q' is required.",
		)
		return
	}

	destinations, err := controller.service.Autocomplete(query)
	if err != nil {
		controller.sendInternalServerError()
		return
	}

	controller.sendSuccess(
		"Destinations retrieved successfully.",
		destinations,
	)
}

func (controller *DestinationController) Nearby() {
	latitude, err := strconv.ParseFloat(
		controller.Ctx.Input.Query("lat"),
		64,
	)
	if err != nil {
		controller.sendBadRequest(
			"Invalid latitude.",
		)
		return
	}

	longitude, err := strconv.ParseFloat(
		controller.Ctx.Input.Query("lon"),
		64,
	)
	if err != nil {
		controller.sendBadRequest(
			"Invalid longitude.",
		)
		return
	}

	radiusKm, err := strconv.ParseFloat(
		controller.Ctx.Input.Query("radius"),
		64,
	)
	if err != nil {
		controller.sendBadRequest(
			"Invalid radius.",
		)
		return
	}

	destinations, err := controller.service.Nearby(
		latitude,
		longitude,
		radiusKm,
	)
	if err != nil {
		controller.sendBadRequest(err.Error())
		return
	}

	controller.sendSuccess(
		"Nearby destinations retrieved successfully.",
		destinations,
	)
}