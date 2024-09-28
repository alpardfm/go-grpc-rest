package rest

import (
	"net/http"
	"strconv"

	"github.com/alpardfm/go-grpc-rest/db"
	"github.com/alpardfm/go-grpc-rest/pb"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func CreateProductHandler(c *gin.Context) {
	var product pb.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.CreateProduct(&db.Product{Name: product.Name, Price: product.Price}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": product})
}

func GetProductHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed convert string to int"})
	}
	product, err := db.GetProduct(int32(idInt))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": product})
}

func UpdateProductHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed convert string to int"})
	}
	var product pb.Product
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.UpdateProduct(&db.Product{ID: int32(idInt), Name: product.Name, Price: product.Price}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": product})
}

func DeleteProductHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed convert string to int"})
	}
	if err := db.DeleteProduct(int32(idInt)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func ListProductsHandler(c *gin.Context) {
	products, err := db.ListProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var productList pb.ProductList
	for _, p := range products {
		productList.Products = append(productList.Products, &pb.Product{
			Id:    p.ID,
			Name:  p.Name,
			Price: p.Price,
		})
	}
	c.JSON(http.StatusOK, productList)
}
