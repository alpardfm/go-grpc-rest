package grpc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/alpardfm/go-grpc-rest/pb"
)

type mockProductServer struct {
	pb.UnimplementedProductServiceServer
}

func (s *mockProductServer) CreateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	return req, nil
}

func (s *mockProductServer) GetProduct(ctx context.Context, req *pb.ProductID) (*pb.Product, error) {
	return &pb.Product{Id: req.Id, Name: "Test Product", Price: 99.99}, nil
}

func (s *mockProductServer) UpdateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	return req, nil
}

func (s *mockProductServer) DeleteProduct(ctx context.Context, req *pb.ProductID) (*pb.ProductID, error) {
	return req, nil
}

func (s *mockProductServer) ListProducts(ctx context.Context, req *emptypb.Empty) (*pb.ProductList, error) {
	return &pb.ProductList{Products: []*pb.Product{
		{Id: 1, Name: "Product 1", Price: 99.99},
		{Id: 2, Name: "Product 2", Price: 79.99},
	}}, nil
}

func TestProductService(t *testing.T) {
	// Setup mock gRPC server
	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, &mockProductServer{})

	// Creating a mock client connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewProductServiceClient(conn)

	// Test CreateProduct
	product := &pb.Product{Name: "Test Product", Price: 99.99}
	resp, err := client.CreateProduct(context.Background(), product)
	require.NoError(t, err)
	assert.Equal(t, product, resp)

	// Test GetProduct
	productID := &pb.ProductID{Id: 3}
	respProduct, err := client.GetProduct(context.Background(), productID)
	require.NoError(t, err)
	assert.Equal(t, &pb.Product{Id: 3, Name: "Test Product", Price: 99.99}, respProduct)

	// Test UpdateProduct
	updatedProduct := &pb.Product{Id: 3, Name: "Updated Product", Price: 89.99}
	resp, err = client.UpdateProduct(context.Background(), updatedProduct)
	require.NoError(t, err)
	assert.Equal(t, updatedProduct, resp)

	// Test DeleteProduct
	deletedProductID, err := client.DeleteProduct(context.Background(), productID)
	require.NoError(t, err)
	assert.Equal(t, productID, deletedProductID)

	// Test ListProducts
	respList, err := client.ListProducts(context.Background(), &emptypb.Empty{})
	require.NoError(t, err)
	assert.Len(t, respList.Products, 2)
	assert.Equal(t, "Product 1", respList.Products[0].Name)
}
