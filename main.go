package main

import (
	"Inshorts/configs"
	"Inshorts/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	//
	configs.ConnectDB()

	routes.UserRoute(e) //add this
	e.Logger.Fatal(e.Start(":6000"))
}
