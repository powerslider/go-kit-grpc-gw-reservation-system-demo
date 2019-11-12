package reservation

import (
	"context"
	gt "github.com/go-kit/kit/transport/grpc"
	errors "github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/error"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/proto"
	"google.golang.org/grpc"
)

func MakeGRPCServer(grpcServer *grpc.Server, s Service) *grpc.Server {
	endpoints := MakeServerEndpoints(s)
	proto.RegisterReservationServiceServer(grpcServer, newGRPCServer(endpoints))
	return grpcServer
}

type grpcServer struct {
	bookReservation                  gt.Handler
	discardReservation               gt.Handler
	//editReservation                  gt.Handler
	getReservationHistoryPerCustomer gt.Handler
}

func (s *grpcServer) BookReservation(ctx context.Context, req *proto.BookReservationRequest) (*proto.BookReservationResponse, error) {
	_, resp, err := s.bookReservation.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.BookReservationResponse), nil
}

func (s *grpcServer) DiscardReservation(ctx context.Context, req *proto.DiscardReservationRequest) (*proto.DiscardReservationResponse, error) {
	_, resp, err := s.discardReservation.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.DiscardReservationResponse), nil
}

//func (s *grpcServer) EditReservation(ctx context.Context, req *proto.EditReservationRequest) (*proto.EditReservationResponse, error) {
//	_, resp, err := s.editReservation.ServeGRPC(ctx, req)
//	if err != nil {
//		return nil, err
//	}
//	return resp.(*proto.EditReservationResponse), nil
//}

func (s *grpcServer) GetReservationHistoryPerCustomer(ctx context.Context, req *proto.GetReservationHistoryPerCustomerRequest) (*proto.GetReservationHistoryPerCustomerResponse, error) {
	_, resp, err := s.getReservationHistoryPerCustomer.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.GetReservationHistoryPerCustomerResponse), nil
}

func newGRPCServer(endpoint Endpoints) proto.ReservationServiceServer {
	return &grpcServer{
		bookReservation: gt.NewServer(
			endpoint.BookReservationEndpoint,
			decodeBookReservationRequest,
			encodeBookReservationResponse,
		),
		discardReservation: gt.NewServer(
			endpoint.DiscardReservationEndpoint,
			decodeDiscardReservationRequest,
			encodeDiscardReservationResponse,
		),
		//editReservation: gt.NewServer(
		//	endpoint.EditReservationEndpoint,
		//	decodeEditReservationRequest,
		//	encodeEditReservationResponse,
		//),
		getReservationHistoryPerCustomer: gt.NewServer(
			endpoint.GetReservationHistoryByCustomerEndpoint,
			decodeGetReservationHistoryPerCustomerRequest,
			encodeGetReservationHistoryPerCustomerResponse,
		),
	}
}

func decodeBookReservationRequest(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*proto.BookReservationRequest)
	return bookReservationRequest{
		CustomerID: int(reply.CustomerId),
		Reservation: reply.Reservation,
	}, nil
}

func decodeDiscardReservationRequest(_ context.Context, grpcReply interface{}) (request interface{}, err error) {
	reply := grpcReply.(*proto.DiscardReservationRequest)
	return discardReservationRequest{
		ReservationID: int(reply.ReservationId),
	}, nil
}

//func decodeEditReservationRequest(_ context.Context, grpcReply interface{}) (request interface{}, err error) {
//	reply := grpcReply.(*proto.EditReservationRequest)
//	return editReservationRequest{
//		ReservationID: int(reply.ReservationId),
//	}, nil
//}

func decodeGetReservationHistoryPerCustomerRequest(_ context.Context, grpcReply interface{}) (request interface{}, err error) {
	reply := grpcReply.(*proto.GetReservationHistoryPerCustomerRequest)
	return getReservationHistoryPerCustomerRequest{
		Limit:  uint(reply.Limit),
		Offset: uint(reply.Offset),
	}, nil
}

func encodeBookReservationResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(bookReservationResponse)
	return &proto.BookReservationResponse{
		Reservation: resp.Reservation,
		Err:         errors.Err2str(resp.Err),
	}, nil
}

func encodeDiscardReservationResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(discardReservationResponse)
	return &proto.DiscardReservationResponse{
		Err: errors.Err2str(resp.Err),
	}, nil
}

//func encodeEditReservationResponse(_ context.Context, response interface{}) (interface{}, error) {
//	resp := response.(editReservationResponse)
//	return &proto.EditReservationResponse{
//		Reservation: resp.Reservation,
//		Err:         errors.Err2str(resp.Err),
//	}, nil
//}

func encodeGetReservationHistoryPerCustomerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(getReservationHistoryPerCustomerResponse)
	return &proto.GetReservationHistoryPerCustomerResponse{
		Reservations: resp.Reservations,
		Err:          errors.Err2str(resp.Err),
	}, nil
}
