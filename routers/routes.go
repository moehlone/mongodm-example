package routers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/moehlone/mongodm_sample/controllers"
	response "github.com/zebresel-com/beego-response"
)

func init() {

	var FilterAuthenticated = func(ctx *context.Context) {

		/**
		 * It is necessary to allow the /login route if you have one.
		 * Until now I did not found a better solution than:
		 *
		 * if strings.Contains(ctx.Request.URL.Path, "/api/login") {
		 *     return
		 *  }
		 *
		 * Fit the comparison to your needs!
		 */

		isAuthorized := true
		response := response.New(ctx)

		if isAuthorized {
			/**
			 * Maybe add a user object here on valid login/authorization for other controllers
			 * e.g. ctx.Input.SetData("user", user)
			 *
			 * In baseController::Prepare you could access the user with:
			 *
			 * if user, ok := self.Ctx.Input.GetData("user").(*models.User); ok {
			 *     self.user = user
			 * }
			 */
		} else {
			/**
			 * This would result in:
			 *
			 * {
			 *  "error": {
			 *      "code": 401,
			 *      "message": "Unauthorized"
			 *      }
			 *  }
			 */
			response.Error(http.StatusUnauthorized)
		}
	}

	// Each route has to pass the filter implemented above
	beego.InsertFilter("/api/*", beego.BeforeRouter, FilterAuthenticated)

	// This is a global namespace for all routes -> http://HOST/api/...
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/user",
			beego.NSRouter("/", &controllers.UserController{}, "get:GetAll"),
			beego.NSRouter("/", &controllers.UserController{}, "post:Create"),
		),
		beego.NSNamespace("/message",
			beego.NSRouter("/", &controllers.MessageController{}, "get:GetAll"),
			beego.NSRouter("/", &controllers.MessageController{}, "post:Create"),
		),
	)
	beego.AddNamespace(ns)
}
