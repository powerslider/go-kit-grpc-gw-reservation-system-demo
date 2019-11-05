package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/powerslider/go-kit-grpc-reservation-system-demo/docs"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/customer"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/storage"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/proto"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// @title Reservation System API
// @version 1.0
// @description Demo service demonstrating Go-Kit.
// @termsOfService http://swagger.io/terms/

// @contact.name Tsvetan Dimitrov
// @contact.email tsvetan.dimitrov23@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {

	var (
		grpcAddr   = flag.String("grpc.addr", ":8080", "gRPC listen address")
		grpcGWAddr = flag.String("grpcgw.addr", ":8081", "gRPC-gateway listen address")
	)
	flag.Parse()

	db, err := storage.NewDB("reservations")
	if err != nil {
		panic(err)
	}

	logger := log.NewLogfmtLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	grpcServer := grpc.NewServer()
	grpcServer = initCustomerGRPCServer(grpcServer, db, logger)
	//grpcServer = initReservationHandler(r, db, logger)

	// The gRPC listener mounts the Go kit gRPC server we created.
	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		os.Exit(1)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "gRPC", "addr", *grpcAddr)
		errs <- grpcServer.Serve(grpcListener)
	}()

	// See https://github.com/grpc/grpc/blob/master/doc/naming.md
	// for gRPC naming standard information.
	dialAddr := fmt.Sprintf("passthrough://localhost/%s", *grpcAddr)
	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, err := grpc.DialContext(
		context.Background(),
		dialAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		logger.Log("Failed to dial server:", err)
	}

	jsonpb := &runtime.JSONPb{
		//EmitDefaults: true,
		Indent:       "  ",
		OrigName:     false,
	}
	mux := http.NewServeMux()

	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
		// This is necessary to get error details properly
		// marshalled in unary requests.
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
	)
	err = proto.RegisterCustomerServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		logger.Log("Failed to register gateway:", err)
	}
	mux.Handle("/", gwmux)
	err = serveSwagger(mux, *grpcGWAddr)
	if err != nil {
		logger.Log("Failed to serve Swagger")
	}

	gwServer := &http.Server{
		Addr:    *grpcGWAddr,
		Handler: mux,
	}

	gwServer.ListenAndServe()

	logger.Log("exit", <-errs)
}

func serveSwagger(mux *http.ServeMux, grpcGWAddr string) error {
	prefix := "/swagger/"
	swaggerURL := fmt.Sprintf("http://localhost%s/swagger/doc.json", grpcGWAddr)
	swaggerHandlerFunc := httpSwagger.Handler(httpSwagger.URL(swaggerURL))
	mux.Handle(prefix, swaggerHandlerFunc)
	return nil
}

func initCustomerGRPCServer(grpcServer *grpc.Server, db *storage.Persistence, logger log.Logger) *grpc.Server {
	r := customer.NewCustomerRepository(*db)
	s := customer.NewCustomerService(r)
	s = customer.LoggingMiddleware(logger)(s)
	return customer.MakeGRPCServer(grpcServer, s)
}

//func initReservationHandler(router *mux.Router, db *storage.Persistence, logger log.Logger) *mux.Router {
//	r := reservation.NewReservationRepository(*db)
//	s := reservation.NewReservationService(r)
//	s = reservation.LoggingMiddleware(logger)(s)
//	return reservation.MakeHTTPHandler(router, s, logger)
//}
