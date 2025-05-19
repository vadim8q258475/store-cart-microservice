package service

import (
	"context"

	productpb "github.com/vadim8q258475/store-product-microservice/gen/v1"
)

type ProductService interface {
	Get(ctx context.Context, productId uint32) (*productpb.Product, error)
}

type productService struct {
	client productpb.ProductServiceClient
}

func NewproductService(client productpb.ProductServiceClient) ProductService {
	return &productService{
		client: client,
	}
}

func (s *productService) Get(ctx context.Context, productId uint32) (*productpb.Product, error) {
	request := &productpb.GetById_Request{Id: productId}
	response, err := s.client.GetById(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Product, nil
}
