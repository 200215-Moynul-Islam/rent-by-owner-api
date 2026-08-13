package controllers

import (
	"net/http"
	"rent-by-owner-api/utils"

	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

func (controller *BaseController) sendSuccess(message string, data interface{}) {
	utils.SendJSONResponse(
		controller.Ctx,
		http.StatusOK,
		true,
		message,
		data,
	)
}

func (controller *BaseController) sendBadRequest(message string) {
	utils.SendJSONResponse(
		controller.Ctx,
		http.StatusBadRequest,
		false,
		message,
		nil,
	)
}

func (controller *BaseController) sendNotFound(message string) {
	utils.SendJSONResponse(
		controller.Ctx,
		http.StatusNotFound,
		false,
		message,
		nil,
	)
}

func (controller *BaseController) sendInternalServerError() {
	utils.SendJSONResponse(
		controller.Ctx,
		http.StatusInternalServerError,
		false,
		"Internal server error.",
		nil,
	)
}