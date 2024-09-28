package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Mengimpor driver PostgreSQL
)

type Product struct {
	ID    int32   `db:"id"`
	Name  string  `db:"name"`
	Price float32 `db:"price"`
}

var DB *sqlx.DB

// InitDB melakukan koneksi ke database PostgreSQL
func InitDB() *sqlx.DB {
	var err error
	DB, err = sqlx.Open("postgres", "user=postgres dbname=db_test password=f3rm0c4rd host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	fmt.Println("Successfully connected to the database!")

	return DB
}

// CreateProduct menambahkan produk baru ke dalam database
func CreateProduct(p *Product) error {
	_, err := DB.NamedExec(`INSERT INTO products (name, price) VALUES (:name, :price)`, p)
	return err
}

// GetProduct mengambil produk berdasarkan ID
func GetProduct(id int32) (*Product, error) {
	var product Product
	err := DB.Get(&product, `SELECT * FROM products WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// UpdateProduct memperbarui produk berdasarkan ID
func UpdateProduct(p *Product) error {
	_, err := DB.NamedExec(`UPDATE products SET name=:name, price=:price WHERE id=:id`, p)
	return err
}

// DeleteProduct menghapus produk berdasarkan ID
func DeleteProduct(id int32) error {
	_, err := DB.Exec(`DELETE FROM products WHERE id=$1`, id)
	return err
}

// ListProducts mengambil semua produk dari database
func ListProducts() ([]Product, error) {
	var products []Product
	err := DB.Select(&products, `SELECT * FROM products`)
	return products, err
}
