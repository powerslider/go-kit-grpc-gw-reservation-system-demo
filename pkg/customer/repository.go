package customer

import (
	"github.com/doug-martin/goqu/v7"
	"github.com/doug-martin/goqu/v7/exec"
	"github.com/jinzhu/copier"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/gen/go/proto"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/apperror"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/storage"
	"time"
)

const (
	defaultLimit  uint = 100
	defaultOffset uint = 0
)

type Entity struct {
	CustomerID  int    `db:"cid" goqu:"skipinsert"`
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	Email       string
	Phone       string
	Created     int64
	LastUpdated int64 `db:"last_updated"`
}

type Repository interface {
	AddCustomer(c *proto.Customer) (*proto.Customer, error)
	RemoveCustomer(cID int) error
	FindAllCustomers(opts *storage.QueryOptions) ([]*proto.Customer, error)
	FindCustomerByID(cID int) (*proto.Customer, error)
}

type customerRepository struct {
	db storage.Persistence
}

func NewCustomerRepository(db storage.Persistence) Repository {
	return &customerRepository{db: db}
}

func (r *customerRepository) AddCustomer(c *proto.Customer) (*proto.Customer, error) {
	created := time.Now().Unix()

	var entity Entity
	err := copier.Copy(&entity, &c)

	result, err := r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		entity.Created = created
		entity.LastUpdated = created
		return tx.From("customer").Insert(entity)
	})
	err = copier.Copy(&c, &entity)
	if err != nil {
		return nil, apperror.Wrap(apperror.DBError, err, "error adding new customer")
	}

	cID, _ := result.LastInsertId()
	c.CustomerId = cID

	return c, nil
}

func (r *customerRepository) RemoveCustomer(cID int) error {
	_, err := r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		return tx.From("customer").Where(goqu.Ex{"cid": cID}).Delete()
	})

	if err != nil {
		return apperror.Wrapf(apperror.DBError, err, "error deleting customers with id %d", cID)
	}
	return nil
}

func (r *customerRepository) FindAllCustomers(opts *storage.QueryOptions) ([]*proto.Customer, error) {
	if opts.Limit == 0 {
		opts.Limit = defaultLimit
	}

	var ee []Entity
	var cc []*proto.Customer
	err := r.db.DB.From("customer").
		Limit(opts.Limit).
		Offset(opts.Offset).
		ScanStructs(&ee)

	for _, entity := range ee {
		var c *proto.Customer
		err = copier.Copy(&c, &entity)
		if err != nil {
			break
		}
		c.CustomerId = int64(entity.CustomerID)
		cc = append(cc, c)
	}

	if err != nil {
		return nil, apperror.Wrapf(apperror.DBError, err, "error getting all customers")
	}
	return cc, nil
}

func (r *customerRepository) FindCustomerByID(cID int) (*proto.Customer, error) {
	var entity Entity
	var c *proto.Customer
	found, err := r.db.DB.From("customer").Where(
		goqu.C("cid").Eq(cID),
	).ScanStruct(&entity)

	if !found {
		return nil, apperror.Newf(apperror.NotFound, "customer with ID %d not found", cID).
			AddContext("CustomerID", "non existent ID")
	}
	err = copier.Copy(&c, &entity)
	if err != nil {
		return nil, apperror.Wrapf(apperror.DBError, err, "error getting customer with ID %d", cID)
	}
	c.CustomerId = int64(entity.CustomerID)

	return c, nil
}
