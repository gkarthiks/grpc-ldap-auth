package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"grpc-ldap-auth/interceptors"
	"grpc-ldap-auth/proto"
	"net"
)

var (
	grpcPort = "6000"
)

type SimpleLDAPService struct {
	proto.UnimplementedSimpleLDAPServiceServer
}

func main() {
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				interceptors.BasicAuthInterceptor,
			),
		),
	)

	logrus.Debugln("registering the CNP Onboard server")
	proto.RegisterSimpleLDAPServiceServer(grpcServer, &SimpleLDAPService{})

	// Start gRPC server
	logrus.Debugf("starting the CNP gRPC server on port: %s\n", grpcPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}
	logrus.Infof("Starting gRPC server on %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Errorf("error serving gRPC: %v", err)
		return
	}

}

func (s *SimpleLDAPService) SayHi(ctx context.Context, request *proto.SayHiRequest) (*proto.SayHiResponse, error) {
	requestString := request.GetMyName()
	var greetTo string
	if len(requestString) > 0 {
		greetTo = requestString
	} else {
		greetTo = "stranger"
	}
	return &proto.SayHiResponse{GreetingResponse: fmt.Sprintf("Hello there %s !", greetTo)}, nil
}
