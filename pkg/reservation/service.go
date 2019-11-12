package reservation

import (
	"context"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/storage"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/proto"
)

type Service interface {
	BookReservation(ctx context.Context, cID int, r *proto.Reservation) (*proto.Reservation, error)
	DiscardReservation(ctx context.Context, rID int) error
	EditReservation(ctx context.Context, rID int, r *proto.Reservation) (*proto.Reservation, error)
	GetReservationHistoryPerCustomer(ctx context.Context, cID int, opts *storage.QueryOptions) ([]proto.Reservation, error)
}

type reservationService struct {
	resRepo Repository
}

func NewReservationService(repo Repository) Service {
	return &reservationService{
		resRepo: repo,
	}
}

func (s *reservationService) BookReservation(ctx context.Context, cID int, r *proto.Reservation) (*proto.Reservation, error) {
	return s.resRepo.AddReservation(cID, r)
}

func (s *reservationService) DiscardReservation(ctx context.Context, rID int) error {
	return s.resRepo.RemoveReservation(rID)
}

func (s *reservationService) EditReservation(ctx context.Context, rID int, res *proto.Reservation) (r *proto.Reservation, err error) {
	return s.resRepo.UpdateReservation(rID, res)
}

func (s *reservationService) GetReservationHistoryPerCustomer(ctx context.Context, cID int, opts *storage.QueryOptions) ([]proto.Reservation, error) {
	return s.resRepo.FindReservationsByCustomerID(cID, opts)
}
