package handler

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/supabase-community/supabase-go"
)

// route terkait lokasi
func RegisterLocationRoutes(e *echo.Echo, client *supabase.Client) {
	e.GET("/api/locations", func(c echo.Context) error {
		var results []Location
		data, _, err := client.From("locations").Select("*", "exact", false).Execute()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal mengambil data lokasi"})
		}
		if err := json.Unmarshal(data, &results); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal parsing data lokasi"})
		}
		return c.JSON(http.StatusOK, results)
	})

	e.GET("/api/locations/:id/tenants", func(c echo.Context) error {
		locationID := c.Param("id")
		var results []Tenant
		data, _, err := client.From("tenants").Select("*", "exact", false).Eq("location_id", locationID).Execute()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal mengambil data tenant"})
		}
		if err := json.Unmarshal(data, &results); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal parsing data tenant"})
		}
		return c.JSON(http.StatusOK, results)
	})
}