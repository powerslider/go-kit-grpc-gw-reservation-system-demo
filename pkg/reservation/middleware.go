package reservation

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/gen/go/proto"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/storage"
	"time"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) BookReservation(ctx context.Context, cID int, r *proto.Reservation) (result *proto.Reservation, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "BookReservation", "id", r.ReservationId, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.BookReservation(ctx, cID, r)
}

func (mw loggingMiddleware) DiscardReservation(ctx context.Context, rID int) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DiscardReservation", "id", rID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DiscardReservation(ctx, rID)
}

func (mw loggingMiddleware) EditReservation(ctx context.Context, rID int, res *proto.Reservation) (r *proto.Reservation, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "EditReservation", "id", rID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.EditReservation(ctx, rID, res)
}

func (mw loggingMiddleware) GetReservationHistoryPerCustomer(ctx context.Context, cID int, opts *storage.QueryOptions) (result []*proto.Reservation, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetReservationHistoryPerCustomer", "id", cID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetReservationHistoryPerCustomer(ctx, cID, opts)
}
