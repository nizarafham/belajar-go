package handler

import (
	"net/http"
	"uper-eats/lib" 

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Untuk pengembangan, bisa diperketat di produksi
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	supabaseClient := lib.NewSupabaseClient()

	// route testing
	e.GET("/api/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	RegisterLocationRoutes(e, supabaseClient)
	RegisterTenantRoutes(e, supabaseClient)
	RegisterMenuRoutes(e, supabaseClient) 
	RegisterAuthRoutes(e, supabaseClient)

	apiV1 := e.Group("/api/v1", lib.AuthMiddleware)
	RegisterOrderRoutes(apiV1, supabaseClient)
	RegisterUserRoutes(apiV1, supabaseClient)

	e.ServeHTTP(w, r)
}