// @APIVersion 1.0.0
// @Title Rent by Owner API
// @Description API for searching and booking rental properties directly from owners.
package routers

import (
	"rent-by-owner-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace(
		"/api",
		beego.NSNamespace(
			"/v1",
			beego.NSRouter(
				"/health",
				&controllers.HealthController{},
				"get:Health",
			),
			beego.NSNamespace(
				"/destinations",
				beego.NSRouter(
					"/search",
					&controllers.DestinationController{},
					"get:Search",
				),
				beego.NSRouter(
					"/autocomplete",
					&controllers.DestinationController{},
					"get:Autocomplete",
				),
				beego.NSRouter(
					"/nearby",
					&controllers.DestinationController{},
					"get:Nearby",
				),
			),
		),
	)
	beego.AddNamespace(ns)
}