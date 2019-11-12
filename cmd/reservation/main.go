package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/powerslider/go-kit-grpc-reservation-system-demo/docs"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/customer"
	"github.com/powerslider/go-kit-grpc-reservation-system-demo/pkg/reservation"
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
	grpcServer = initReservationGRPCServer(grpcServer, db, logger)

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
	backgroundCtx := context.Background()
	conn, err := grpc.DialContext(
		backgroundCtx,
		dialAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		logger.Log("Failed to dial server:", err)
	}

	jsonpb := &runtime.JSONPb{
		//EmitDefaults: true,
		Indent:   "  ",
		OrigName: false,
	}
	//runtime.HTTPError = CustomHTTPError
	mux := http.NewServeMux()

	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
		// This is necessary to get error details properly
		// marshalled in unary requests.
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
	)
	err = proto.RegisterCustomerServiceHandler(backgroundCtx, gwmux, conn)
	err = proto.RegisterReservationServiceHandler(backgroundCtx, gwmux, conn)
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

func initReservationGRPCServer(grpcServer *grpc.Server, db *storage.Persistence, logger log.Logger) *grpc.Server {
	r := reservation.NewReservationRepository(*db)
	s := reservation.NewReservationService(r)
	s = reservation.LoggingMiddleware(logger)(s)
	return reservation.MakeGRPCServer(grpcServer, s)
}

//type errorBody struct {
//	Err string `json:"error,omitempty"`
//}
//
//func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
//	const fallback = `{"error": "failed to marshal error message"}`
//
//	w.Header().Set("Content-type", marshaler.ContentType())
//	w.WriteHeader(runtime.HTTPStatusFromCode(status.Code(err)))
//	jErr := json.NewEncoder(w).Encode(errorBody{
//		Err: status.Convert(err).Message(),
//	})
//
//	if jErr != nil {
//		w.Write([]byte(fallback))
//	}
//}
