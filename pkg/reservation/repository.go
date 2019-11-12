package reservation

import (
	"github.com/doug-martin/goqu/v7"
	"github.com/doug-martin/goqu/v7/exec"
	"github.com/jinzhu/copier"
	errors "github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/error"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/storage"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/proto"
	"time"
)

type Entity struct {
	ReservationID   int    `db:"rid" goqu:"skipinsert"`
	SeatCount       int    `db:"seat_count"`
	StartTime       string `db:"start_time"`
	ReservationName string `db:"reservation_name"`
	CustomerID      int    `db:"customer_id"`
	Phone           string
	Comments        string
	Created         int64
	LastUpdated     int64  `db:"last_updated"`
}

type Repository interface {
	AddReservation(cID int, r *proto.Reservation) (*proto.Reservation, error)
	RemoveReservation(rID int) error
	UpdateReservation(rID int, r *proto.Reservation) (*proto.Reservation, error)
	FindReservationsByCustomerID(cID int, opts *storage.QueryOptions) ([]proto.Reservation, error)
}

type reservationRepository struct {
	db storage.Persistence
}

func NewReservationRepository(db storage.Persistence) Repository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) AddReservation(cID int, res *proto.Reservation) (*proto.Reservation, error) {
	created := time.Now().Unix()

	var entity Entity
	err := copier.Copy(&entity, &res)

	result, err := r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		entity.Created = created
		entity.LastUpdated = created
		entity.CustomerID = cID
		return tx.From("reservation").Insert(entity)
	})
	err = copier.Copy(&res, &entity)
	if err != nil {
		return nil, errors.DBError.Wrap(err, "error adding new reservation")
	}

	rID, _ := result.LastInsertId()
	res.ReservationId = int64(rID)

	return res, nil
}

func (r *reservationRepository) RemoveReservation(rID int) error {
	_, err := r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		return tx.From("reservation").Where(goqu.Ex{"rid": rID}).Delete()
	})

	if err != nil {
		return errors.DBError.Wrapf(err, "error deleting reservation with ID %d", rID)
	}
	return nil
}

func (r *reservationRepository) UpdateReservation(rID int, res *proto.Reservation) (result *proto.Reservation, err error) {
	lastUpdated := time.Now().Unix()

	var entity Entity
	err = copier.Copy(&entity, &res)

	_, err = r.db.Tx(func(tx *goqu.TxDatabase) exec.QueryExecutor {
		entity.LastUpdated = lastUpdated
		return tx.From("reservation").Update(entity)
	})

	_, err = r.db.DB.From("reservation").Where(
		goqu.C("rid").Eq(rID),
	).ScanStruct(&entity)

	err = copier.Copy(&result, &entity)
	if err != nil {
		return result, errors.DBError.Wrapf(err, "error updating reservation with ID %d", rID)
	}

	return result, nil
}

func (r *reservationRepository) FindReservationsByCustomerID(cID int, opts *storage.QueryOptions) (rr []proto.Reservation, err error) {
	var ee []Entity

	err = r.db.DB.From("reservation").
		Select("reservation.*").
		Join(
			goqu.T("customer"),
			goqu.On(goqu.Ex{
				"reservation.customer_id": goqu.I("customer.cid"),
			})).
		Order(goqu.C("last_updated").Desc()).
		Limit(opts.Limit).
		Offset(opts.Offset).
		ScanStructs(&ee)

	for _, entity := range ee {
		var r proto.Reservation
		err = copier.Copy(&r, &entity)
		if err != nil {
			break
		}
		r.ReservationId = int64(entity.ReservationID) 
		rr = append(rr, r)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "error fetching reservations for customer with ID %d", cID)
	}
	return rr, nil
}
