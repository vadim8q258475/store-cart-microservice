package grpc

import (
	"context"

	gen "github.com/vadim8q258475/store-cart-microservice/gen/v1"
	productpb "github.com/vadim8q258475/store-product-microservice/gen/v1"
	userpb "github.com/vadim8q258475/store-user-microservice/gen/v1"
)

type UserService interface {
	Get(ctx context.Context, userId uint32) (*userpb.User, error)
}
type CartService interface {
	Create(ctx context.Context, userId uint32) (uint32, error)
	Delete(ctx context.Context, cartId uint32) error
	Get(ctx context.Context, cartID uint32) (*gen.Cart, error)
	GetByUserId(ctx context.Context, userId uint32) (*gen.Cart, error)
	List(ctx context.Context) ([]*gen.Cart, error)
	Add(ctx context.Context, cartId, productId, qty uint32) error
	Remove(ctx context.Context, cartId, productId, qty uint32) error
}
type ProductService interface {
	Get(ctx context.Context, productId uint32) (*productpb.Product, error)
}

type GrpcService struct {
	gen.UnimplementedCartServiceServer
	cartService CartService
}

func NewGrpcService(CartService CartService) *GrpcService {
	return &GrpcService{
		cartService: CartService,
	}
}

func (g *GrpcService) Add(ctx context.Context, request *gen.Add_Request) (*gen.Add_Response, error) {
	err := g.cartService.Add(ctx, request.CartId, request.ProductId, request.Qty)
	if err != nil {
		return nil, err
	}
	return &gen.Add_Response{CartId: request.CartId}, nil
}
func (g *GrpcService) Create(ctx context.Context, request *gen.Create_Request) (*gen.Create_Response, error) {
	id, err := g.cartService.Create(ctx, request.UserId)
	if err != nil {
		return nil, err
	}
	return &gen.Create_Response{CartId: id}, err
}
func (g *GrpcService) Delete(ctx context.Context, request *gen.Delete_Request) (*gen.Delete_Response, error) {
	err := g.cartService.Delete(ctx, request.CartId)
	if err != nil {
		return nil, err
	}
	return &gen.Delete_Response{Success: true}, nil
}
func (g *GrpcService) Get(ctx context.Context, request *gen.Get_Request) (*gen.Get_Response, error) {
	cart, err := g.cartService.Get(ctx, request.CartId)
	if err != nil {
		return nil, err
	}
	return &gen.Get_Response{Cart: cart}, nil
}
func (g *GrpcService) GetByUserId(ctx context.Context, request *gen.GetByUserId_Request) (*gen.GetByUserId_Response, error) {
	cart, err := g.cartService.GetByUserId(ctx, request.UserId)
	if err != nil {
		return nil, err
	}
	return &gen.GetByUserId_Response{Cart: cart}, nil
}
func (g *GrpcService) List(ctx context.Context, request *gen.List_Request) (*gen.List_Response, error) {
	carts, err := g.cartService.List(ctx)
	if err != nil {
		return nil, err
	}
	return &gen.List_Response{Carts: carts}, nil
}
func (g *GrpcService) Remove(ctx context.Context, request *gen.Remove_Request) (*gen.Remove_Response, error) {
	err := g.cartService.Remove(ctx, request.CartId, request.ProductId, request.Qty)
	if err != nil {
		return nil, err
	}
	return &gen.Remove_Response{CartId: request.CartId}, nil
}
