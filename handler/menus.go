package handler

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/supabase-community/supabase-go"
)

// route terkait menu
func RegisterMenuRoutes(e *echo.Echo, client *supabase.Client) {
	e.GET("/api/tenants/:id/menus", getMenusByTenant(client))

	e.GET("/api/menus/:id", getMenuDetail(client))
}

func getMenusByTenant(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		tenantID := c.Param("id")
		var results []Menu
		
		// Query ke tabel 'menus' dengan filter berdasarkan tenant_id
		data, _, err := client.From("menus").Select("*", "exact", false).Eq("tenant_id", tenantID).Execute()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal mengambil data menu"})
		}

		if err := json.Unmarshal(data, &results); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal parsing data menu"})
		}

		return c.JSON(http.StatusOK, results)
	}
}

// getMenuDetail menangani permintaan untuk mendapatkan detail satu menu spesifik.
func getMenuDetail(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		menuID := c.Param("id")
		var result Menu

		// Query ke tabel 'menus' dengan filter berdasarkan id menu
		data, _, err := client.From("menus").Select("*", "exact", false).Eq("id", menuID).Single().Execute()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Menu tidak ditemukan"})
		}

		if err := json.Unmarshal(data, &result); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal parsing data menu"})
		}

		return c.JSON(http.StatusOK, result)
	}
}