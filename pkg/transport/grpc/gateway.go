package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpswagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
	"time"
)

type GatewayServer struct {
	Name           string
	GatewayAddr    string
	GRPCServerAddr string
	serverInstance *http.Server
}

func NewGatewayServer(name string, gwAddr string, grpcServerAddr string, serviceServerRegistrarFunc func(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux)) *GatewayServer {
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				Indent:    "  ",
				Multiline: true, // Optional, implied by presence of "Indent".
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		// This is necessary to get apperrors details properly
		// marshalled in unary requests.
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
	)

	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	ctx := context.Background()
	dialAddr := fmt.Sprintf("dns:///%s", grpcServerAddr)
	conn, err := grpc.DialContext(
		ctx,
		dialAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)

	if err != nil {
		panic(err)
	}
	serviceServerRegistrarFunc(ctx, conn, gwmux)

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	mux.Handle("/swagger/", swaggerHandler(gwAddr))

	return &GatewayServer{
		Name:           name,
		GatewayAddr:    gwAddr,
		GRPCServerAddr: grpcServerAddr,
		serverInstance: &http.Server{
			Addr:    gwAddr,
			Handler: mux,
		},
	}
}

func (s *GatewayServer) Start() {
	go func() {
		if err := s.serverInstance.ListenAndServe(); err != nil {
			log.Printf("gRPC-gateway server closing with message: %v", err)
		}
	}()
	log.Printf("[Start] %s gRPC-gateway server on port %s started\n", s.Name, s.GatewayAddr)
}

func (s *GatewayServer) Shutdown() {
	log.Printf("[Shutdown] %s gRPC-gateway server is shutting down\n", s.Name)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.serverInstance.Shutdown(ctx); err != nil {
		log.Fatalf("Unable to stop gRPC-gateway server: %v", err)
	}
}

func swaggerHandler(port string) http.HandlerFunc {
	swaggerURL := fmt.Sprintf("http://localhost%s/swagger/doc.json", port)
	return httpswagger.Handler(httpswagger.URL(swaggerURL))
}
