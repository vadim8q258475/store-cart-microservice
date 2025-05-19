package repo

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("cart not found")

type CartRepo interface {
	Create(ctx context.Context, cart Cart) (uint32, error)
	Delete(ctx context.Context, id uint32) error
	Get(ctx context.Context, id uint32) (Cart, error)
	List(ctx context.Context) ([]Cart, error)
}

type cartRepo struct {
	db *sqlx.DB
}

func NewCartRepo(db *sqlx.DB) CartRepo {
	return &cartRepo{
		db: db,
	}
}

func (r *cartRepo) Create(ctx context.Context, cart Cart) (uint32, error) {
	var id uint32
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO carts (user_id)
		VALUES ($1) RETURNING id`,
		cart.UserId,
	).Scan(&id)
	return id, err
}
func (r *cartRepo) Delete(ctx context.Context, id uint32) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM carts WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return err
}
func (r *cartRepo) Get(ctx context.Context, id uint32) (Cart, error) {
	var cart Cart
	err := r.db.GetContext(ctx, &cart, `SELECT id, user_id FROM carts WHERE id = $1`, id)
	return cart, err
}
func (r *cartRepo) List(ctx context.Context) ([]Cart, error) {
	var carts []Cart
	err := r.db.SelectContext(ctx, &carts, `SELECT id, user_id FROM carts`)
	return carts, err
}
