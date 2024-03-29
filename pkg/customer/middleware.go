package customer

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

func (mw loggingMiddleware) RegisterCustomer(ctx context.Context, c *proto.Customer) (result *proto.Customer, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "RegisterCustomer", "id", c.CustomerId, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.RegisterCustomer(ctx, c)
}

func (mw loggingMiddleware) UnregisterCustomer(ctx context.Context, cID int) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "UnregisterCustomer", "id", cID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.UnregisterCustomer(ctx, cID)
}

func (mw loggingMiddleware) GetAllCustomers(ctx context.Context, opts *storage.QueryOptions) (result []*proto.Customer, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAllCustomers", "limit", opts.Limit, "offset", opts.Offset, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetAllCustomers(ctx, opts)
}

func (mw loggingMiddleware) GetCustomerByID(ctx context.Context, cID int) (result *proto.Customer, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetCustomerByID", "id", cID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetCustomerByID(ctx, cID)
}
