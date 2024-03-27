package coredogcontroller

import (
	"log"
	"net"

	"github.com/DomineCore/coredog/pb"
	"google.golang.org/grpc"
)

const (
	Port = "8443"
)

func ListenAndServe() {
	server := grpc.NewServer()
	pb.RegisterCoreFileServiceServer(server, &CorefileService{})

	lis, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("listen on %s error:%v", Port, err)
	}
	server.Serve(lis)
}
