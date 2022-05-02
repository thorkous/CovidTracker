package routes

import (
	"Inshorts/controllers"
	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	//All routes related to users comes here
	e.GET("/covid/:state", controllers.GetStateActiveCases) //add this
	e.GET("/covid/:x/:y", controllers.GetStateActiveCasesUsingPosition)

}
