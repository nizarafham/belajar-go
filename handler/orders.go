package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"uper-eats/lib" 

	"github.com/labstack/echo/v4"
	"github.com/supabase-community/supabase-go"
	"github.com/xendit/xendit-go/v3/invoice"
)

func RegisterOrderRoutes(g *echo.Group, client *supabase.Client) {
	g.POST("/orders", createOrder(client))
	g.GET("/orders", getOrderHistory(client))
	g.GET("/orders/:id", getOrderDetail(client))
}

type CreateOrderRequest struct {
	TenantID        int64                  `json:"tenant_id"`
	OrderType       string                 `json:"order_type"` 
	DeliveryAddress string                 `json:"delivery_address,omitempty"`
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

		// 1. Ambil data tenant
		var tenantData struct {
			Name string `json:"name"`
		}
		dataTenant, _, err := client.From("tenants").Select("name", "exact", false).Eq("id", fmt.Sprintf("%d", req.TenantID)).Single().Execute()
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Tenant tidak ditemukan."})
		}
		json.Unmarshal(dataTenant, &tenantData)

		// 2. Ambil data user
		var user User
		userData, _, err := client.From("users").Select("email", "exact", false).Eq("id", fmt.Sprintf("%d", userID)).Single().Execute()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal mendapatkan data pengguna"})
		}
		json.Unmarshal(userData, &user)

		// 3. Hitung total harga
		var totalPrice float64
		var orderItemsToInsert []OrderItem
		for _, itemReq := range req.Items {
			var menu Menu
			data, _, err := client.From("menus").Select("price, is_available, name", "exact", false).Eq("id", fmt.Sprintf("%d", itemReq.MenuID)).Single().Execute()
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

		// 4. Insert ke tabel 'orders'
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
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal membuat pesanan di database"})
		}
		json.Unmarshal(orderData, &insertedOrder)
		orderID := insertedOrder[0].ID

		// 5. Insert item
		for i := range orderItemsToInsert {
			orderItemsToInsert[i].OrderID = orderID
		}
		_, _, err = client.From("order_items").Insert(orderItemsToInsert, false, "error", "", "exact").Execute()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal menyimpan item pesanan"})
		}

		// 6. Buat Invoice (tanpa ForUserId / Fees jika tidak didukung SDK kamu)
		desc := fmt.Sprintf("Pemesanan di %s - Order #%d", tenantData.Name, orderID)
		createInvoiceRequest := invoice.CreateInvoiceRequest{
			ExternalId:  strconv.FormatInt(orderID, 10),
			Amount:      float32(totalPrice),
			PayerEmail:  &user.Email,
			Description: &desc,
		}

		invoiceRequest := lib.XenditClient.InvoiceApi.
			CreateInvoice(c.Request().Context()).
			CreateInvoiceRequest(createInvoiceRequest) 

		resp, _, err := invoiceRequest.Execute()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Gagal membuat invoice Xendit: " + err.Error(),
			})
		}


		// 7. Response ke frontend
		return c.JSON(http.StatusCreated, echo.Map{
			"message":     "Pesanan berhasil dibuat, silakan lakukan pembayaran.",
			"order_id":    orderID,
			"total_price": totalPrice,
			"payment_url": resp.InvoiceUrl,
		})
	}
}


func getOrderHistory(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("userID").(int64)

		var results []Order
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
