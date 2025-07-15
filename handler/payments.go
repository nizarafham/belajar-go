package handler

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/supabase-community/supabase-go"
)

func RegisterPaymentRoutes(e *echo.Echo, client *supabase.Client) {
	e.POST("/api/payments/xendit-notification", handleXenditNotification(client))
}

func handleXenditNotification(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1. Verifikasi Webhook (Sangat Penting untuk Keamanan)
		callbackToken := c.Request().Header.Get("X-CALLBACK-TOKEN")
		expectedToken := os.Getenv("XENDIT_WEBHOOK_VERIFICATION_TOKEN")

		if callbackToken != expectedToken {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Verifikasi webhook gagal"})
		}

		// 2. Proses Notifikasi
		var notificationPayload map[string]interface{}
		if err := c.Bind(&notificationPayload); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Request tidak valid"})
		}

		orderID, ok := notificationPayload["external_id"].(string)
		if !ok {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "external_id tidak ditemukan"})
		}

		status, _ := notificationPayload["status"].(string)
		var newStatus string

		if status == "PAID" || status == "SETTLED" {
			newStatus = "paid" 
		} else if status == "EXPIRED" {
			newStatus = "cancelled" 
		}

		// 3. Update Status Pesanan di Database
		if newStatus != "" {
			_, _, err := client.From("orders").
				Update(map[string]interface{}{"status": newStatus, "updated_at": "now()"}, "", "exact").
				Eq("id", orderID).
				Execute()

			if err != nil {
				c.Logger().Errorf("Gagal update status order %s: %v", orderID, err)
				// Tetap kirim 200 agar Xendit tidak coba lagi
			}
		}

		return c.JSON(http.StatusOK, echo.Map{"status": "notifikasi diterima"})
	}
}