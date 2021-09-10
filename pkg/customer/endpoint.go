package customer

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/gen/go/proto"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/storage"
)

type Endpoints struct {
	RegisterCustomerEndpoint   endpoint.Endpoint
	UnregisterCustomerEndpoint endpoint.Endpoint
	GetAllCustomersEndpoint    endpoint.Endpoint
	GetCustomerByIDEndpoint    endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		RegisterCustomerEndpoint:   MakeRegisterCustomerEndpoint(s),
		UnregisterCustomerEndpoint: MakeUnregisterCustomerEndpoint(s),
		GetAllCustomersEndpoint:    MakeGetAllCustomersEndpoint(s),
		GetCustomerByIDEndpoint:    MakeGetCustomerByIDEndpoint(s),
	}
}

type unregisterCustomerRequest struct {
	CustomerID int
}

type unregisterCustomerResponse struct {
	Err error `json:"err,omitempty"`
}

func (r unregisterCustomerResponse) Failed() error { return r.Err }

func MakeUnregisterCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(unregisterCustomerRequest)
		e := s.UnregisterCustomer(ctx, req.CustomerID)
		return unregisterCustomerResponse{
			Err: e,
		}, nil
	}
}

type registerCustomerRequest struct {
	Customer *proto.Customer
}

type registerCustomerResponse struct {
	Customer *proto.Customer `json:"customer,omitempty"`
	Err      error           `json:"err,omitempty"`
}

func (r registerCustomerResponse) Failed() error { return r.Err }

func MakeRegisterCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerCustomerRequest)
		c, e := s.RegisterCustomer(ctx, req.Customer)
		return registerCustomerResponse{
			Customer: c,
			Err:      e,
		}, nil
	}
}

type getAllCustomersRequest struct {
	Limit  uint
	Offset uint
}

type getAllCustomersResponse struct {
	Customers []*proto.Customer `json:"customers,omitempty"`
	Err       error            `json:"err,omitempty"`
}

func (r getAllCustomersResponse) Failed() error { return r.Err }

func MakeGetAllCustomersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getAllCustomersRequest)
		cc, e := s.GetAllCustomers(ctx, &storage.QueryOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		})
		return getAllCustomersResponse{
			Customers: cc,
			Err:       e,
		}, nil
	}
}

type getCustomerByIDRequest struct {
	CustomerID int
}

type getCustomerByIDResponse struct {
	Customer *proto.Customer `json:"customer,omitempty"`
	Err      error           `json:"err,omitempty"`
}

func (r getCustomerByIDResponse) Failed() error { return r.Err }

func MakeGetCustomerByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getCustomerByIDRequest)
		c, e := s.GetCustomerByID(ctx, req.CustomerID)
		return getCustomerByIDResponse{
			Customer: c,
			Err:      e,
		}, nil
	}
}
