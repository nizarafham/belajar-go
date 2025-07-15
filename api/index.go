package handler

import (
	"net/http"

	"uper-eats/lib"
	"uper-eats/handler" 

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	e := echo.New()

	// Middleware global
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // for development, change to specific origins in production
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	supabaseClient := lib.NewSupabaseClient()

	// Route testing
	e.GET("/api/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// Public routes
	handler.RegisterLocationRoutes(e, supabaseClient)
	handler.RegisterTenantRoutes(e, supabaseClient)
	handler.RegisterMenuRoutes(e, supabaseClient)
	handler.RegisterAuthRoutes(e, supabaseClient)

	// Protected routes (dengan JWT Auth)
	apiV1 := e.Group("/api/v1", lib.AuthMiddleware)
	handler.RegisterOrderRoutes(apiV1, supabaseClient)
	handler.RegisterUserRoutes(apiV1, supabaseClient)

	e.ServeHTTP(w, r)
}
