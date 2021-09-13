package main

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/powerslider/go-kit-grpc-reservation-system-demo/docs"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/gen/go/proto"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/customer"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/reservation"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/storage"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/transport"
	grpctransport "github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/transport/grpc"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/transport/tcp"
	"google.golang.org/grpc"
	"os"
)

func main() {
	var (
		grpcGwAddr     = ":8081"
		grpcServerAddr = "0.0.0.0:8080"
		appName        = "grpc-gw-reservation-system-server"
		appGwName      = "grpc-gw-reservation-system-gateway"
	)

	db, err := storage.NewDB("reservations")
	if err != nil {
		panic(err)
	}

	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	grpcTCPListener, err := tcp.NewTCPListener(grpcServerAddr)
	if err != nil {
		panic(err)
	}

	grpcServer := grpctransport.NewGRPCServer(appName, grpcServerAddr, grpcTCPListener, func(s *grpc.Server) *grpc.Server {
		s = initCustomerGRPCServer(s, db, logger)
		s = initReservationGRPCServer(s, db, logger)
		return s
	})
	grpcServer.Start()

	grpcGwServer := grpctransport.NewGatewayServer(appGwName, grpcGwAddr, grpcServerAddr, func(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
		err = proto.RegisterCustomerServiceHandler(ctx, mux, conn)
		err = proto.RegisterReservationServiceHandler(ctx, mux, conn)
		if err != nil {
			panic(err)
		}
	})
	grpcGwServer.Start()

	transport.WaitForShutdownSignal()

	grpcGwServer.Shutdown()
	grpcServer.Shutdown()
}

func initCustomerGRPCServer(grpcServer *grpc.Server, db *storage.Persistence, logger log.Logger) *grpc.Server {
	r := customer.NewCustomerRepository(*db)
	s := customer.NewCustomerService(r)
	s = customer.LoggingMiddleware(logger)(s)
	return customer.MakeGRPCServer(grpcServer, s)
}

func initReservationGRPCServer(grpcServer *grpc.Server, db *storage.Persistence, logger log.Logger) *grpc.Server {
	r := reservation.NewReservationRepository(*db)
	s := reservation.NewReservationService(r)
	s = reservation.LoggingMiddleware(logger)(s)
	return reservation.MakeGRPCServer(grpcServer, s)
}
