package repo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type CartProductRepo interface {
	Create(ctx context.Context, cartProduct CartProduct) (uint32, error)
	Update(ctx context.Context, cartProduct CartProduct) (uint32, error)
	GetByCartId(ctx context.Context, cartId uint32) ([]CartProduct, error)
	AddToCart(ctx context.Context, cartId, cartProductId uint32) error
	Delete(ctx context.Context, cartProductId uint32) error
	GetByProductId(ctx context.Context, productId, cartId uint32) (CartProduct, error)
}

type cartProductRepo struct {
	db *sqlx.DB
}

func NewCartProductRepo(db *sqlx.DB) CartProductRepo {
	return &cartProductRepo{
		db: db,
	}
}

func (r *cartProductRepo) Create(ctx context.Context, cartProduct CartProduct) (uint32, error) {
	var id uint32
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO cart_products (product_id, qty)
		VALUES ($1, $2) RETURNING id`,
		cartProduct.ProductId, cartProduct.Qty,
	).Scan(&id)
	return id, err
}
func (r *cartProductRepo) AddToCart(ctx context.Context, cartId, cartProductId uint32) error {
	var id uint32
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO cart_cart_products (cart_id, cart_product_id)
		VALUES ($1, $2) RETURNING id`,
		cartId, cartProductId,
	).Scan(&id)
	return err
}
func (r *cartProductRepo) Update(ctx context.Context, cartProduct CartProduct) (uint32, error) {
	var id uint32
	err := r.db.QueryRowContext(ctx,
		`UPDATE cart_products
		SET qty = $1
		WHERE id = $2 
		RETURNING id`,
		cartProduct.Qty, cartProduct.Id,
	).Scan(&id)
	return id, err
}
func (r *cartProductRepo) Delete(ctx context.Context, cartProductId uint32) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM cart_cart_products WHERE cart_product_id = $1`, cartProductId)
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
	result, err = r.db.ExecContext(ctx, `DELETE FROM cart_products WHERE id = $1`, cartProductId)
	if err != nil {
		return err
	}
	rows, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
func (r *cartProductRepo) GetByProductId(ctx context.Context, productId, cartId uint32) (CartProduct, error) {
	cartProducts, err := r.GetByCartId(ctx, cartId)
	if err != nil {
		return CartProduct{}, err
	}
	for _, cartProduct := range cartProducts {
		if cartProduct.ProductId == productId {
			return cartProduct, nil
		}
	}
	return CartProduct{}, status.Error(codes.NotFound, "cart product does not exist")
}
func (r *cartProductRepo) GetByCartId(ctx context.Context, cartId uint32) ([]CartProduct, error) {
	var cartProductIds []uint32
	err := r.db.SelectContext(ctx, &cartProductIds,
		`SELECT cart_product_id FROM cart_cart_products WHERE cart_id = $1`, cartId)
	if err != nil {
		return nil, err
	}
	fmt.Println(cartProductIds)
	cartProducts := make([]CartProduct, len(cartProductIds))
	for i, id := range cartProductIds {
		err := r.db.GetContext(ctx, &cartProducts[i],
			`SELECT id, product_id, qty FROM cart_products WHERE id = $1`, id)
		if err != nil {
			return nil, err
		}
	}
	return cartProducts, nil
}
