package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"uper-eats/lib"

	"github.com/labstack/echo/v4"
	"github.com/supabase-community/supabase-go"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAuthRoutes(e *echo.Echo, client *supabase.Client) {
	e.POST("/api/auth/register", registerUser(client))
	e.POST("/api/auth/login", loginUser(client))
}

func RegisterUserRoutes(g *echo.Group, client *supabase.Client) {
	g.GET("/users/me", getProfile(client))
}

func registerUser(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		type RegisterRequest struct {
			FullName    string `json:"full_name"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number"`
			Password    string `json:"password"`
			Role        string `json:"role,omitempty"`
		}

		req := new(RegisterRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Request tidak valid"})
		}

		// Validasi input wajib
		if req.FullName == "" || req.Email == "" || req.PhoneNumber == "" || req.Password == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Semua field wajib diisi"})
		}

		if req.Role == "" {
			req.Role = RoleUser
		}
		if !IsValidRole(req.Role) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "Role tidak valid. Pilihan: user, tenant_owner, driver",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal memproses password"})
		}

		newUser := User{
			FullName:     req.FullName,
			Email:        req.Email,
			PhoneNumber:  req.PhoneNumber,
			PasswordHash: string(hashedPassword),
			Role:         req.Role,
			CreatedAt:    time.Now(),
		}

		var results []User
		data, _, err := client.From("users").Insert(newUser, false, "error", "", "exact").Execute()
		if err != nil {
			fmt.Println("Insert error:", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Gagal mendaftarkan pengguna, mungkin email atau nomor sudah terdaftar",
			})
		}

		json.Unmarshal(data, &results)
		return c.JSON(http.StatusCreated, results[0])
	}
}

func loginUser(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		type LoginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		req := new(LoginRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Request tidak valid"})
		}

		var result User
		data, _, err := client.From("users").Select("*", "exact", false).Eq("email", req.Email).Single().Execute()
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Email atau password salah"})
		}
		json.Unmarshal(data, &result)

		err = bcrypt.CompareHashAndPassword([]byte(result.PasswordHash), []byte(req.Password))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Email atau password salah"})
		}

		token, err := lib.GenerateJWT(result.ID, result.Role)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal membuat token"})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Login berhasil",
			"token":   token,
			"user": echo.Map{
				"id":        result.ID,
				"full_name": result.FullName,
				"email":     result.Email,
				"role":      result.Role,
			},
		})
	}
}

func getProfile(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Get("userID").(int64)

		var result User
		data, _, err := client.From("users").Select("id,full_name,email,phone_number,role,created_at", "exact", false).Eq("id", fmt.Sprintf("%d", userID)).Single().Execute()
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Pengguna tidak ditemukan"})
		}
		json.Unmarshal(data, &result)

		return c.JSON(http.StatusOK, result)
	}
}