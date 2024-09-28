package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	rest "github.com/alpardfm/go-grpc-rest/api/rest" // Ganti dengan import path yang sesuai

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Product struct {
	ID    int32   `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/products", rest.CreateProductHandler)
	router.GET("/products/:id", rest.GetProductHandler)
	return router
}

func TestCreateProductHandler(t *testing.T) {
	// Mocking database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	rest.DB = sqlxDB // Set global DB di package rest

	// Mocking product data
	product := &Product{Name: "Test Product", Price: 99.99}

	// Setting expectation for SQL query
	mock.ExpectExec(`INSERT INTO products (name, price) VALUES ($1, $2)`).
		WithArgs(product.Name, product.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Creating request with body
	productJSON, err := json.Marshal(product)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(productJSON))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Setting up the response recorder and router
	rr := httptest.NewRecorder()
	router := setupRouter() // Make sure /products route is correctly initialized
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusCreated, rr.Code) // Use http.StatusCreated (201) if you're creating a new resource

	// Check if all expectations for the mock DB were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductHandler(t *testing.T) {
	// Mocking database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	rest.DB = sqlxDB // Set global DB di package rest

	product := &Product{ID: 1, Name: "Test Product", Price: 99.99}

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(product.ID, product.Name, product.Price)

	mock.ExpectQuery(`SELECT \* FROM products WHERE id=\$1`).
		WithArgs(product.ID).
		WillReturnRows(rows)

	// Creating request
	req, err := http.NewRequest("GET", "/products/1", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
