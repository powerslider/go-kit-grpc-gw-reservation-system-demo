package customer

import (
	"context"
	gt "github.com/go-kit/kit/transport/grpc"
	errors "github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/error"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/proto"
)

type grpcServer struct {
	registerCustomer   gt.Handler
	unregisterCustomer gt.Handler
	getAllCustomers    gt.Handler
	getCustomerByID    gt.Handler
}

func (s *grpcServer) RegisterCustomer(context.Context, *proto.RegisterCustomerRequest) (*proto.RegisterCustomerResponse, error) {
	panic("implement me")
}

func (s *grpcServer) UnregisterCustomer(context.Context, *proto.UnregisterCustomerRequest) (*proto.UnregisterCustomerResponse, error) {
	panic("implement me")
}

func (s *grpcServer) GetAllCustomers(context.Context, *proto.GetAllCustomersRequest) (*proto.GetAllCustomersResponse, error) {
	panic("implement me")
}

func (s *grpcServer) GetCustomerByID(context.Context, *proto.GetCustomerByIDRequest) (*proto.GetCustomerByIDResponse, error) {
	panic("implement me")
}

func NewGRPCServer(_ context.Context, endpoint Endpoints) proto.CustomerServiceServer {
	return &grpcServer{
		registerCustomer: gt.NewServer(
			endpoint.RegisterCustomerEndpoint,
			decodeRegisterCustomerRequest,
			encodeRegisterCustomerResponse,
		),
		unregisterCustomer: gt.NewServer(
			endpoint.UnregisterCustomerEndpoint,
			decodeUnregisterCustomerRequest,
			encodeUnregisterCustomerResponse,
		),
		getCustomerByID: gt.NewServer(
			endpoint.GetCustomerByIDEndpoint,
			decodeGetCustomerByIDRequest,
			encodeGetCustomerByIDResponse,
		),
		getAllCustomers: gt.NewServer(
			endpoint.GetAllCustomersEndpoint,
			decodeGetAllCustomersRequest,
			encodeGetAllCustomersResponse,
		),
	}
}

func decodeRegisterCustomerRequest(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.RegisterCustomerRequest)
	return registerCustomerRequest{
		Customer: reply.Customer,
	}, nil
}

func decodeUnregisterCustomerRequest(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.UnregisterCustomerRequest)
	return unregisterCustomerRequest{
		CustomerID: int(reply.CustomerId),
	}, nil
}

func decodeGetCustomerByIDRequest(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.GetCustomerByIDRequest)
	return getCustomerByIDRequest{
		CustomerID: int(reply.CustomerId),
	}, nil
}

func decodeGetAllCustomersRequest(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.GetAllCustomersRequest)
	return getAllCustomersRequest{
		Limit:  uint(reply.Limit),
		Offset: uint(reply.Offset),
	}, nil
}

func encodeRegisterCustomerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(registerCustomerResponse)
	return &proto.RegisterCustomerResponse{
		Customer: resp.Customer,
		Err:      errors.Err2str(resp.Err),
	}, nil
}

func encodeUnregisterCustomerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(unregisterCustomerResponse)
	return &proto.UnregisterCustomerResponse{
		Err: errors.Err2str(resp.Err),
	}, nil
}

func encodeGetCustomerByIDResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(getCustomerByIDResponse)
	return &proto.GetCustomerByIDResponse{
		Customer: &resp.Customer,
		Err:      errors.Err2str(resp.Err),
	}, nil
}

func encodeGetAllCustomersResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(getAllCustomersResponse)
	return &proto.GetAllCustomersResponse{
		Customers: resp.Customers,
		Err:       errors.Err2str(resp.Err),
	}, nil
}
