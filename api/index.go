package handler

import (
	"fmt"
	"net/http"

	. "github.com/tbxark/g4vercel"
)

// Middleware untuk mengatur header CORS
func corsMiddleware() HandlerFunc {
	return func(c *Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Jika method OPTIONS, cukup return 200 tanpa lanjut ke handler
		if c.Method == http.MethodOptions {
			c.Writer.WriteHeader(http.StatusOK)
			return
		}

		c.Next()
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	server := New()

	// Tambahkan middleware CORS
	server.Use(corsMiddleware())

	// Recovery middleware (biar kalau error tetap kirim JSON)
	server.Use(Recovery(func(err interface{}, c *Context) {
		if httpError, ok := err.(HttpError); ok {
			c.JSON(httpError.Status, H{
				"message": httpError.Error(),
			})
		} else {
			message := fmt.Sprintf("%s", err)
			c.JSON(500, H{
				"message": message,
			})
		}
	}))

	server.GET("/", func(context *Context) {
		context.JSON(200, H{
			"message": "OK",
		})
	})

	server.GET("/hello", func(context *Context) {
		name := context.Query("name")
		if name == "" {
			context.JSON(400, H{
				"message": "name not found",
			})
		} else {
			context.JSON(200, H{
				"data": fmt.Sprintf("Hello %s!", name),
			})
		}
	})

	server.GET("/user/:id", func(context *Context) {
		context.JSON(400, H{
			"data": H{
				"id": context.Param("id"),
			},
		})
	})

	server.GET("/long/long/long/path/*test", func(context *Context) {
		context.JSON(200, H{
			"data": H{
				"url": context.Path,
			},
		})
	})

	server.Handle(w, r)
}
