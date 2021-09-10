package customer

import (
	"context"
	gt "github.com/go-kit/kit/transport/grpc"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/gen/go/proto"
	errors "github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/apperror"
	"google.golang.org/grpc"
)

func MakeGRPCServer(grpcServer *grpc.Server, s Service) *grpc.Server {
	endpoints := MakeServerEndpoints(s)
	proto.RegisterCustomerServiceServer(grpcServer, newGRPCServer(endpoints))
	return grpcServer
}

type grpcServer struct {
	registerCustomer   gt.Handler
	unregisterCustomer gt.Handler
	getAllCustomers    gt.Handler
	getCustomerByID    gt.Handler
}

func (s *grpcServer) RegisterCustomer(ctx context.Context, req *proto.RegisterCustomerRequest) (*proto.RegisterCustomerResponse, error) {
	_, resp, err := s.registerCustomer.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.RegisterCustomerResponse), nil
}

func (s *grpcServer) UnregisterCustomer(ctx context.Context, req *proto.UnregisterCustomerRequest) (*proto.UnregisterCustomerResponse, error) {
	_, resp, err := s.unregisterCustomer.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.UnregisterCustomerResponse), nil
}

func (s *grpcServer) GetAllCustomers(ctx context.Context, req *proto.GetAllCustomersRequest) (*proto.GetAllCustomersResponse, error) {
	_, resp, err := s.getAllCustomers.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.GetAllCustomersResponse), nil
}

func (s *grpcServer) GetCustomerByID(ctx context.Context, req *proto.GetCustomerByIDRequest) (*proto.GetCustomerByIDResponse, error) {
	_, resp, err := s.getCustomerByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.GetCustomerByIDResponse), nil
}

func newGRPCServer(endpoint Endpoints) proto.CustomerServiceServer {
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
		Customer: resp.Customer,
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
