package policyservice

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func (ps *PolicyServer) grpcServe(l net.Listener) error {
	creds, err := credentials.NewServerTLSFromFile(ps.Config.CertFile, ps.Config.KeyFile)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	reflection.Register(grpcServer)
	RegisterPolicyServerServer(grpcServer, ps)
	return grpcServer.Serve(l)
}
