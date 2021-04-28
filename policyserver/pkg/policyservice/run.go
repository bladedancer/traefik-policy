package policyservice

import (
	"fmt"
	"net"
)

type PolicyServer struct {
	Config *PolicyServiceConfig
}

func (ps *PolicyServer) Run() error {
	bind := fmt.Sprintf("%s:%d", ps.Config.Host, ps.Config.Port)
	listener, err := net.Listen("tcp4", bind)
	if err != nil {
		log.Fatal(err)
	}

	// m := cmux.New(listener)
	// grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	// httpListener := m.Match(cmux.Any() /*cmux.HTTP1Fast()*/)

	// g := new(errgroup.Group)
	// g.Go(func() error { return ps.grpcServe(grpcListener) })
	// g.Go(func() error { return ps.httpServe(httpListener) })
	// g.Go(func() error { return m.Serve() })

	// log.Infof("Listening on %s", bind)
	// return g.Wait()

	log.Infof("Listening on %s", bind)
	return ps.grpcServe(listener)
}
