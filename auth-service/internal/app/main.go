package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mfajri11/family-catering-micro-service/auth-service/config"
	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/core/service"
	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/handler/rpc"
	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/handler/rpc/pb"
	"github.com/mfajri11/family-catering-micro-service/auth-service/pkg/server"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run() error {
	ctx, cancelF := context.WithCancel(context.Background())
	defer cancelF()
	// prepare connection for grpc and http server (TODO: create connection pooling for grpc)
	grpcClientConn, grpcServerConn := prepareGrpcConn(ctx)
	eg, ctx := errgroup.WithContext(ctx)

	// runtime mux as a gateway which convert http request to rpc call
	rmux := runtime.NewServeMux(runtime.WithForwardResponseOption(rpc.WriteSIDToCookieIFSuccess))
	// register to default mux (maybe using chi in order to utilize middleware)
	// mux := http.NewServeMux()
	// mux.Handle("/", rmux)

	// TODO: crate server struct which implement graceful shutdown and seamless restart
	httpSrv := server.New(config.Cfg.App.Address(), rmux)

	// construct grpc server and service (the rpc handler struct)
	grpcSrv := grpc.NewServer()
	svc := rpc.NewAuthHandler(&service.AuthService{})

	// client for grpc call purpose
	client := pb.NewAuthClient(grpcClientConn)

	// register service and server
	pb.RegisterAuthServer(grpcSrv, svc)
	// register client to be used in handler
	err := pb.RegisterAuthHandlerClient(ctx, rmux, client)
	if err != nil {
		panic(err)
	}
	// only allow 2 goroutines for this group
	eg.SetLimit(2)

	// run grpc server and http server independently
	eg.Go(func() error {
		return grpcSrv.Serve(grpcServerConn)
	})

	eg.Go(func() error {
		// connection for http server will be created within Serve function
		return httpSrv.Serve()
	})

	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)

	// SIGTERM and INTERRUPT/SIGINT for graceful shutdown
	// SIGHUP for hot reload (option via endpoint and db call, binary through another yaml file and env when calling the binary)
	signal.Notify(sigChan, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	errChan <- eg.Wait()

	select {
	case <-sigChan:
		grpcSrv.GracefulStop()
		return httpSrv.Shutdown()
	case err := <-errChan:
		return fmt.Errorf("app.Run: %w", err)
	}

}

func prepareGrpcConn(ctx context.Context) (grpcClientConn *grpc.ClientConn, grpcServerConn net.Listener) {

	// TODO: use grpc pooling
	// construct the connection
	grpcClientConn, err := grpc.DialContext(ctx, ":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	// defer grpcClientConn.Close()
	grpcServerConn, err = net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	// defer grpcServerConn.Close()

	return grpcClientConn, grpcServerConn
}
