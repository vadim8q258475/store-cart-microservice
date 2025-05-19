package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vadim8q258475/store-cart-microservice/config"
)

var cartProductSchema = `
CREATE TABLE IF NOT EXISTS cart_products (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    qty INTEGER NOT NULL DEFAULT 0 CHECK (qty >= 0)
)`

var cartSchema = `
CREATE TABLE IF NOT EXISTS carts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE  
)`

var cartCartProducts = `
CREATE TABLE IF NOT EXISTS cart_cart_products (
    id SERIAL PRIMARY KEY,
    cart_id INTEGER NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    cart_product_id INTEGER NOT NULL REFERENCES cart_products(id) ON DELETE CASCADE,
    UNIQUE(cart_id, cart_product_id)
)`

func createTables(db *sqlx.DB) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(cartProductSchema); err != nil {
		return err
	}

	if _, err := tx.Exec(cartSchema); err != nil {
		return err
	}

	if _, err := tx.Exec(cartCartProducts); err != nil {
		return err
	}

	return tx.Commit()
}

func InitDB(cfg config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = createTables(db)
	return db, err
}
