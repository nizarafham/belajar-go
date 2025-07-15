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
			FullName string `json:"full_name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		req := new(RegisterRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Request tidak valid"})
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal memproses password"})
		}

		newUser := User{
			FullName:     req.FullName,
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
			Role:         "student", 
			CreatedAt:    time.Now(),
		}

		var results []User
		data, _, err := client.From("users").Insert(newUser, false, "error", "", "exact").Execute()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal mendaftarkan pengguna, mungkin email sudah terdaftar"})
		}

		json.Unmarshal(data, &results)
		return c.JSON(http.StatusCreated, results[0])
	}
}

// Handler untuk login
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

		// Bandingkan password
		err = bcrypt.CompareHashAndPassword([]byte(result.PasswordHash), []byte(req.Password))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Email atau password salah"})
		}

		// Buat JWT
		token, err := lib.GenerateJWT(result.ID, result.Role)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Gagal membuat token"})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Login berhasil",
			"token":   token,
		})
	}
}

// Handler untuk mendapatkan profil pengguna
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