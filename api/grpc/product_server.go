package grpc

import (
	"context"

	"github.com/alpardfm/go-grpc-rest/db"
	"github.com/alpardfm/go-grpc-rest/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductServer struct {
	pb.UnimplementedProductServiceServer
}

func (s *ProductServer) CreateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	product := db.Product{Name: req.Name, Price: req.Price}
	if err := db.CreateProduct(&product); err != nil {
		return nil, err
	}

	return req, nil
}

func (s *ProductServer) GetProduct(ctx context.Context, req *pb.ProductID) (*pb.Product, error) {
	product, err := db.GetProduct(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Product{Id: product.ID, Name: product.Name, Price: product.Price}, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	product := db.Product{ID: req.Id, Name: req.Name, Price: req.Price}
	if err := db.UpdateProduct(&product); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *ProductServer) DeleteProduct(ctx context.Context, req *pb.ProductID) (*pb.ProductID, error) {
	if err := db.DeleteProduct(req.Id); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *ProductServer) ListProducts(ctx context.Context, req *emptypb.Empty) (*pb.ProductList, error) {
	products, err := db.ListProducts()
	if err != nil {
		return nil, err
	}
	var productList pb.ProductList
	for _, p := range products {
		productList.Products = append(productList.Products, &pb.Product{
			Id:    p.ID,
			Name:  p.Name,
			Price: p.Price,
		})
	}
	return &productList, nil
}
