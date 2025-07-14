package main

import (
	"net/http"

	"github.com/gin-contrib/cors" // Import CORS middleware
	"github.com/gin-gonic/gin"
)

func main() {
	// Membuat router Gin dengan konfigurasi default
	router := gin.Default()

	// Gunakan CORS Middleware
	// cors.Default() akan mengizinkan semua koneksi selama development.
	// Ini adalah perbaikan untuk masalah loading di Chrome.
	router.Use(cors.Default())

	// Membuat endpoint sederhana untuk testing
	router.GET("/ping", func(c *gin.Context) {
		// Server akan merespon dengan JSON
		c.JSON(http.StatusOK, gin.H{
			"message": "pong! Halo dari server Go (Lokal)!",
		})
	})

	// Menjalankan server di port 8080
	router.Run(":8080")
}
