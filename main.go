package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Handler
func callback(c echo.Context) error {
	var v map[string]interface{}
	m := map[string]interface{}{}
	m["request"] = c.Request().RequestURI
	err := c.Bind(&v)
	m["error"] = err
	m["payload"] = v
	return c.JSON(http.StatusOK, m)
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/oauth/callback", callback)

	// Start server
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
