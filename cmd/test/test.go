package main

import (
	"context"

	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	dsn := "host=localhost user=postgres password=example dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	println(err)

	ctx := context.Background()

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	err = gorm.G[Product](db).Create(ctx, &Product{Code: "D42", Price: 100})

	// Read
	product, err := gorm.G[Product](db).Where("id = ?", 1).First(ctx) // find product with integer primary key

	fmt.Println(product)

	// Update - update product's price to 200
	num, err := gorm.G[Product](db).Where("id = ?", product.ID).Update(ctx, "Price", 200)

	println(num)
	// Update - update multiple fields

	// Delete - delete product

}
