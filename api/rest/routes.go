package rest

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan semua route REST API
func RegisterRoutes(router *gin.Engine) {
	// Tambahkan route untuk CRUD produk
	router.POST("/products", CreateProductHandler)       // Tambah produk baru
	router.GET("/products/:id", GetProductHandler)       // Ambil produk berdasarkan ID
	router.PUT("/products/:id", UpdateProductHandler)    // Update produk berdasarkan ID
	router.DELETE("/products/:id", DeleteProductHandler) // Hapus produk berdasarkan ID
	router.GET("/products", ListProductsHandler)         // Ambil semua produk
}
