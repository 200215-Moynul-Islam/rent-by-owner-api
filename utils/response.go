package utils

import (
	beegoCtx "github.com/beego/beego/v2/server/web/context"
)

// APIResponse represents the standard API response format.
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    interface{}    `json:"data,omitempty"`
}

// SendJSONResponse sends a standardized JSON response.
func SendJSONResponse(ctx *beegoCtx.Context, status int, success bool, message string, data interface{}) {
	ctx.Output.SetStatus(status)

	response := APIResponse{
		Success: success,
		Message: message,
		Data:    data,
	}

	ctx.Output.JSON(response, false, false)
}