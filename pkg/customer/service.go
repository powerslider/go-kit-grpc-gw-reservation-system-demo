package customer

import (
	"context"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/storage"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/proto"
)

type Service interface {
	RegisterCustomer(ctx context.Context, c *proto.Customer) (*proto.Customer, error)
	UnregisterCustomer(ctx context.Context, cID int) error
	GetAllCustomers(ctx context.Context, opts *storage.QueryOptions) ([]proto.Customer, error)
	GetCustomerByID(ctx context.Context, cID int) (proto.Customer, error)
}

//type Customer struct {
//	CustomerID  int    `json:"customerId" db:"cid" goqu:"skipinsert"`
//	FirstName   string `json:"firstName" db:"first_name"`
//	LastName    string `json:"lastName" db:"last_name"`
//	Email       string `json:"email"`
//	Phone       string `json:"phone"`
//	Created     int64  `json:"created"`
//	LastUpdated int64  `json:"lastUpdated" db:"last_updated"`
//}

type customerService struct {
	custRepo Repository
}

func NewCustomerService(repo Repository) Service {
	return &customerService{
		custRepo: repo,
	}
}

func (s *customerService) RegisterCustomer(ctx context.Context, c *proto.Customer) (*proto.Customer, error) {
	return s.custRepo.AddCustomer(c)
}

func (s *customerService) UnregisterCustomer(ctx context.Context, cID int) error {
	return s.custRepo.RemoveCustomer(cID)
}

func (s *customerService) GetAllCustomers(ctx context.Context, opts *storage.QueryOptions) ([]proto.Customer, error) {
	return s.custRepo.FindAllCustomers(opts)
}

func (s *customerService) GetCustomerByID(ctx context.Context, cID int) (proto.Customer, error) {
	return s.custRepo.FindCustomerByID(cID)
}
