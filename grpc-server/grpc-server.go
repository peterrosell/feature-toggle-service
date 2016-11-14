package main

import (
	"net"
	"github.com/golang/glog"

	api "github.com/peterrosell/feature-toggle-service/api-impl"
	"google.golang.org/grpc"
)


func Run() error {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	api.RegisterFeatureToggleService(s)

	s.Serve(l)
	return nil
}

func main() {
 defer glog.Flush()

 if err := Run(); err != nil {
   glog.Fatal(err)
 }
}
