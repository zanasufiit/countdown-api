package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Students struct {
	Count int `json:"count"`
  }

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.Any("/*", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Students { Count: 250 })
	})

	e.Logger.Fatal(e.Start(":" + port()))
}

func port() (port string) {
	port = "80"

	val, present := os.LookupEnv("PORT")
	if present {
		tPort, err := strconv.Atoi(val)
		if err == nil {
			port = strconv.Itoa(tPort)
		}
	}

	return
}
