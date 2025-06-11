package service

import (
	"context"
	"fmt"

	gen "github.com/vadim8q258475/store-cart-microservice/gen/v1"
	"github.com/vadim8q258475/store-cart-microservice/iternal/repo"
	productService "github.com/vadim8q258475/store-cart-microservice/iternal/service/product"
	userService "github.com/vadim8q258475/store-cart-microservice/iternal/service/user"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type CartService interface {
	Create(ctx context.Context, userId uint32) (uint32, error)
	Delete(ctx context.Context, cartId uint32) error
	Get(ctx context.Context, cartID uint32) (*gen.Cart, error)
	GetByUserId(ctx context.Context, userId uint32) (*gen.Cart, error)
	List(ctx context.Context) ([]*gen.Cart, error)
	Add(ctx context.Context, cartId, productId, qty uint32) error
	Remove(ctx context.Context, cartId, productId, qty uint32) error
}

type cartService struct {
	productService  productService.ProductService
	userService     userService.UserService
	cartRepo        repo.CartRepo
	cartProductRepo repo.CartProductRepo
}

func NewCartService(productService productService.ProductService,
	userService userService.UserService,
	cartRepo repo.CartRepo,
	cartProductRepo repo.CartProductRepo) CartService {
	return &cartService{
		productService:  productService,
		userService:     userService,
		cartRepo:        cartRepo,
		cartProductRepo: cartProductRepo,
	}
}

func (s *cartService) Create(ctx context.Context, userId uint32) (uint32, error) {
	_, err := s.userService.Get(ctx, userId)
	if err != nil {
		return uint32(0), err
	}
	cart := repo.Cart{UserId: userId}
	return s.cartRepo.Create(ctx, cart)
}
func (s *cartService) Delete(ctx context.Context, cartId uint32) error {
	return s.cartRepo.Delete(ctx, cartId)
}
func (s *cartService) GetCartProducts(ctx context.Context, cartId uint32) ([]*gen.CartProduct, uint32, error) {
	cartProductModels, err := s.cartProductRepo.GetByCartId(ctx, cartId)
	var total uint32
	if err != nil {
		return nil, total, err
	}
	cartProducts := make([]*gen.CartProduct, len(cartProductModels))
	fmt.Println(cartProductModels)
	for i, cartProductModel := range cartProductModels {
		product, err := s.productService.Get(ctx, cartProductModel.ProductId)
		if err != nil {
			return nil, total, err
		}
		total += uint32(product.Price) * cartProductModel.Qty
		cartProducts[i] = &gen.CartProduct{
			Id: cartProductModel.Id,
			Product: &gen.Product{
				Id:          product.Id,
				Name:        product.Name,
				Description: product.Description,
				Qty:         product.Qty,
				Price:       product.Price,
				Category: &gen.Category{
					Id:          product.Category.Id,
					Name:        product.Category.Name,
					Description: product.Category.Description,
				},
			},
			Qty: cartProductModel.Qty,
		}
	}
	return cartProducts, total, nil
}
func (s *cartService) Get(ctx context.Context, cartID uint32) (*gen.Cart, error) {
	cartModel, err := s.cartRepo.Get(ctx, cartID)
	if err != nil {
		return nil, err
	}
	cartProducts, total, err := s.GetCartProducts(ctx, cartID)
	if err != nil {
		return nil, err
	}
	cart := &gen.Cart{
		Id:       cartModel.Id,
		UserId:   cartModel.UserId,
		Products: cartProducts,
		Total:    total,
	}
	return cart, err
}
func (s *cartService) GetByUserId(ctx context.Context, userID uint32) (*gen.Cart, error) {
	cartModel, err := s.cartRepo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}
	cartProducts, total, err := s.GetCartProducts(ctx, cartModel.Id)
	if err != nil {
		return nil, err
	}
	cart := &gen.Cart{
		Id:       cartModel.Id,
		UserId:   cartModel.UserId,
		Products: cartProducts,
		Total:    total,
	}
	return cart, err
}
func (s *cartService) List(ctx context.Context) ([]*gen.Cart, error) {
	cartModels, err := s.cartRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	carts := make([]*gen.Cart, len(cartModels))
	for i, cartModel := range cartModels {
		cartProducts, total, err := s.GetCartProducts(ctx, cartModel.Id)
		if err != nil {
			return nil, err
		}
		carts[i] = &gen.Cart{
			Id:       cartModel.Id,
			UserId:   cartModel.UserId,
			Products: cartProducts,
			Total:    total,
		}
	}
	return carts, nil
}
func (s *cartService) Add(ctx context.Context, cartId, productId, qty uint32) error {
	cart, err := s.cartRepo.Get(ctx, cartId)
	if err != nil {
		return err
	}
	product, err := s.productService.Get(ctx, productId)
	if err != nil {
		return err
	}
	if product.Qty < int32(qty) {
		return status.Error(codes.InvalidArgument, "product.qty < qty")
	}
	cartProduct, getCartProductErr := s.cartProductRepo.GetByProductId(ctx, productId, cart.Id)
	if getCartProductErr != nil {
		st, ok := status.FromError(getCartProductErr)
		if !ok {
			return getCartProductErr
		}
		code := st.Code()
		if code != codes.NotFound {
			return getCartProductErr
		}
	}

	if cartProduct.Qty+qty > uint32(product.Qty) {
		return status.Error(codes.InvalidArgument, "cartProduct.Qty + qty > product.Qty")
	}

	cartProduct.Qty += qty

	if getCartProductErr != nil {
		cartProduct.ProductId = productId
		cartProductId, err := s.cartProductRepo.Create(ctx, cartProduct)
		if err != nil {
			return err
		}
		err = s.cartProductRepo.AddToCart(ctx, cart.Id, cartProductId)
		if err != nil {
			return err
		}
	} else {
		_, err = s.cartProductRepo.Update(ctx, cartProduct)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *cartService) Remove(ctx context.Context, cartId, productId, qty uint32) error {
	cart, err := s.cartRepo.Get(ctx, cartId)
	if err != nil {
		return err
	}
	product, err := s.productService.Get(ctx, productId)
	if err != nil {
		return err
	}
	cartProduct, err := s.cartProductRepo.GetByProductId(ctx, product.Id, cart.Id)
	if err != nil {
		return err
	}
	if cartProduct.Qty > qty {
		cartProduct.Qty -= qty
		_, err = s.cartProductRepo.Update(ctx, cartProduct)
		if err != nil {
			return err
		}
		return nil
	}
	return s.cartProductRepo.Delete(ctx, cartProduct.Id)
}
