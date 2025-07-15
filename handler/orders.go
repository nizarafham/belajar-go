package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/supabase-community/supabase-go"
)

// route terkait order
func RegisterOrderRoutes(g *echo.Group, client *supabase.Client) {
	g.POST("/orders", createOrder(client))
	g.GET("/orders", getOrderHistory(client))
	g.GET("/orders/:id", getOrderDetail(client))
}

// Struct untuk request pembuatan order
type CreateOrderRequest struct {
	TenantID        int64                     `json:"tenant_id"`
	OrderType       string                    `json:"order_type"` // 'pickup' or 'delivery'
	DeliveryAddress string                    `json:"delivery_address,omitempty"`
	Items           []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	MenuID   int64 `json:"menu_id"`
	Quantity int   `json:"quantity"`
}

func createOrder(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("userID").(int64)
		req := new(CreateOrderRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Request tidak valid"})
		}

		// 1. Validasi dan hitung total harga (SERVER-SIDE)
		var totalPrice float64 = 0
		var orderItemsToInsert []OrderItem

		for _, itemReq := range req.Items {
			var menu Menu
			data, _, err := client.From("menus").Select("price, is_available", "exact", false).Eq("id", fmt.Sprintf("%d", itemReq.MenuID)).Single().Execute()
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("Menu dengan ID %d tidak ditemukan", itemReq.MenuID)})
			}
			json.Unmarshal(data, &menu)

			if !menu.IsAvailable {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("Menu %s sedang tidak tersedia", menu.Name)})
			}

			itemPrice := menu.Price * float64(itemReq.Quantity)
			totalPrice += itemPrice

			orderItemsToInsert = append(orderItemsToInsert, OrderItem{
				MenuID:       itemReq.MenuID,
				Quantity:     itemReq.Quantity,
				PricePerItem: menu.Price,
			})
		}

		// 2. Buat record di tabel 'orders'
		newOrder := Order{
			UserID:          userID,
			TenantID:        req.TenantID,
			TotalPrice:      totalPrice,
			OrderType:       req.OrderType,
			Status:          "pending_payment",
			DeliveryAddress: req.DeliveryAddress,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

		var insertedOrder []Order
		orderData, _, err := client.From("orders").Insert(newOrder, false, "error", "", "exact").Execute()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal membuat pesanan"})
		}
		json.Unmarshal(orderData, &insertedOrder)
		orderID := insertedOrder[0].ID

		// 3. Buat record di tabel 'order_items'
		for i := range orderItemsToInsert {
			orderItemsToInsert[i].OrderID = orderID
		}
		_, _, err = client.From("order_items").Insert(orderItemsToInsert, false, "error", "", "exact").Execute()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal menyimpan detail pesanan"})
		}

		// 4. (SIMULASI) Buat transaksi Midtrans dan kembalikan URL pembayaran
		// Di production harus memanggil SDK Midtrans di sini
		paymentRedirectURL := fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v2/vtweb/%d-dummy-token", orderID)

		return c.JSON(http.StatusCreated, echo.Map{
			"message":            "Pesanan berhasil dibuat, silakan lakukan pembayaran.",
			"order_id":           orderID,
			"total_price":        totalPrice,
			"payment_redirect_url": paymentRedirectURL,
		})
	}
}

func getOrderHistory(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("userID").(int64)
		
		var results []Order
		// Mengambil data order beserta info tenant
		data, _, err := client.From("orders").Select("*, tenant_id(*)", "exact", false).Eq("user_id", fmt.Sprintf("%d", userID)).Execute()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal mengambil riwayat pesanan"})
		}

		if err := json.Unmarshal(data, &results); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal parsing data pesanan"})
		}

		return c.JSON(http.StatusOK, results)
	}
}

func getOrderDetail(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("userID").(int64)
		orderID := c.Param("id")

		var result Order
		// Query kompleks untuk mengambil order, item-itemnya, dan info menu & tenant
		query := "*, tenant_id(*), order_items(*, menu_id(*))"
		data, _, err := client.From("orders").Select(query, "exact", false).Eq("id", orderID).Eq("user_id", fmt.Sprintf("%d", userID)).Single().Execute()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Pesanan tidak ditemukan atau bukan milik Anda"})
		}

		if err := json.Unmarshal(data, &result); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal parsing detail pesanan"})
		}

		return c.JSON(http.StatusOK, result)
	}
}