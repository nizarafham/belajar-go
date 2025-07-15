package handler

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/supabase-community/supabase-go"
)

// route tenant
func RegisterTenantRoutes(e *echo.Echo, client *supabase.Client) {
	e.GET("/api/tenants/:id", func(c echo.Context) error {
		tenantID := c.Param("id")
		var result Tenant
		data, _, err := client.From("tenants").Select("*", "exact", false).Eq("id", tenantID).Single().Execute()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Tenant tidak ditemukan"})
		}
		if err := json.Unmarshal(data, &result); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal parsing data tenant"})
		}
		return c.JSON(http.StatusOK, result)
	})

	e.GET("/api/tenants/:id/menus", func(c echo.Context) error {
		tenantID := c.Param("id")
		var results []Menu
		data, _, err := client.From("menus").Select("*", "exact", false).Eq("tenant_id", tenantID).Execute()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal mengambil data menu"})
		}
		if err := json.Unmarshal(data, &results); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal parsing data menu"})
		}
		return c.JSON(http.StatusOK, results)
	})
}