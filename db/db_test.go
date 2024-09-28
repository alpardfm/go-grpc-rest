package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	DB = sqlxDB // Set global DB di package db

	product := &Product{Name: "Test Product", Price: 99.99}

	mock.ExpectExec(`INSERT INTO products \(name, price\) VALUES \(\$1, \$2\)`).
		WithArgs(product.Name, product.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = CreateProduct(product)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	DB = sqlxDB // Set global DB di package db

	product := &Product{ID: 1, Name: "Test Product", Price: 99.99}

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(product.ID, product.Name, product.Price)

	mock.ExpectQuery(`SELECT \* FROM products WHERE id=\$1`).
		WithArgs(product.ID).
		WillReturnRows(rows)

	result, err := GetProduct(product.ID)
	require.NoError(t, err)
	assert.Equal(t, product, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	DB = sqlxDB // Set global DB di package db

	product := &Product{ID: 1, Name: "Updated Product", Price: 89.99}

	mock.ExpectExec(`UPDATE products SET name=\$1, price=\$2 WHERE id=\$3`).
		WithArgs(product.Name, product.Price, product.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = UpdateProduct(product)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	DB = sqlxDB // Set global DB di package db

	productID := int32(1)

	mock.ExpectExec(`DELETE FROM products WHERE id=\$1`).
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = DeleteProduct(productID)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	DB = sqlxDB // Set global DB di package db

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Product 1", 99.99).
		AddRow(2, "Product 2", 79.99)

	mock.ExpectQuery(`SELECT \* FROM products`).
		WillReturnRows(rows)

	products, err := ListProducts()
	require.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, "Product 2", products[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}
