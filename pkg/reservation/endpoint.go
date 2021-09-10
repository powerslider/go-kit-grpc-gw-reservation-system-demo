package reservation

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/gen/go/proto"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/storage"
)

type Endpoints struct {
	BookReservationEndpoint                 endpoint.Endpoint
	DiscardReservationEndpoint              endpoint.Endpoint
	EditReservationEndpoint                 endpoint.Endpoint
	GetReservationHistoryByCustomerEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		BookReservationEndpoint:                 MakeBookReservationEndpoint(s),
		DiscardReservationEndpoint:              MakeDiscardReservationEndpoint(s),
		EditReservationEndpoint:                 MakeEditReservationEndpoint(s),
		GetReservationHistoryByCustomerEndpoint: MakeGetReservationHistoryPerCustomerEndpoint(s),
	}
}

type discardReservationRequest struct {
	ReservationID int
}

type discardReservationResponse struct {
	Err error `json:"err,omitempty"`
}

func (r discardReservationResponse) Failed() error { return r.Err }

func MakeDiscardReservationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(discardReservationRequest)
		e := s.DiscardReservation(ctx, req.ReservationID)
		return discardReservationResponse{
			Err: e,
		}, nil
	}
}

type bookReservationRequest struct {
	CustomerID  int
	Reservation *proto.Reservation
}

type bookReservationResponse struct {
	Reservation *proto.Reservation `json:"reservation,omitempty"`
	Err         error              `json:"err,omitempty"`
}

func (r bookReservationResponse) Failed() error { return r.Err }

func MakeBookReservationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(bookReservationRequest)
		r, e := s.BookReservation(ctx, req.CustomerID, req.Reservation)
		return bookReservationResponse{
			Reservation: r,
			Err:         e,
		}, nil
	}
}

type getReservationHistoryPerCustomerRequest struct {
	CustomerID int
	Limit      uint
	Offset     uint
}

type getReservationHistoryPerCustomerResponse struct {
	Reservations []*proto.Reservation `json:"reservations,omitempty"`
	Err          error               `json:"err,omitempty"`
}

func (r getReservationHistoryPerCustomerResponse) Failed() error { return r.Err }

func MakeGetReservationHistoryPerCustomerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getReservationHistoryPerCustomerRequest)
		rr, e := s.GetReservationHistoryPerCustomer(ctx, req.CustomerID, &storage.QueryOptions{
			Limit:  req.Limit,
			Offset: req.Offset,
		})
		return getReservationHistoryPerCustomerResponse{
			Reservations: rr,
			Err:          e,
		}, nil
	}
}

type editReservationRequest struct {
	ReservationID int
	Reservation   *proto.Reservation
}

type editReservationResponse struct {
	Reservation *proto.Reservation `json:"reservation"`
	Err         error              `json:"err,omitempty"`
}

func (r editReservationResponse) Failed() error { return r.Err }

func MakeEditReservationEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(editReservationRequest)
		r, e := s.EditReservation(ctx, req.ReservationID, req.Reservation)
		return editReservationResponse{
			Reservation: r,
			Err:         e,
		}, nil
	}
}
