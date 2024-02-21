package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	gecko "github.com/superoo7/go-gecko/v3"
)

func main() {
	e := echo.New()

	// Templ handler for the homepage
	e.GET("/", func(c echo.Context) error {
		templ.Handler(homepage()).ServeHTTP(c.Response().Writer, c.Request())

		return nil
	})

	e.GET("/price", func(c echo.Context) error {
		ticker := c.QueryParam("ticker")
		if ticker == "" {
			return c.String(http.StatusBadRequest, "Ticker query parameter is required")
		}

		cg := gecko.NewClient(nil)
		price, err := cg.SimplePrice([]string{ticker}, []string{"usd"})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		priceUSD, exists := (*price)[ticker]["usd"]
		if !exists {
			return c.String(http.StatusNotFound, "Price for the specified ticker not found")
		}
		priceUSDStr := fmt.Sprintf("%f", priceUSD)
		templ.Handler(displayPrice(priceUSDStr)).ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	e.Logger.Fatal(e.Start(":8080"))
}
