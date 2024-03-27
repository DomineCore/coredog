package coredogwatcher

import (
	"fmt"
	"log"
	"os"

	"github.com/DomineCore/coredog/pb"

	"google.golang.org/grpc"
)

const (
	ControllerPort = "8443"
	ControllerHost = "localhost"
)

func getControllerPort() string {
	return os.Getenv("CONTROLLER_PORT")
}

func getControllerHost() string {
	return os.Getenv("CONTROLLER_HOST")
}

func NewCoreFileServiceClient() (pb.CoreFileServiceClient, *grpc.ClientConn) {
	controllerPort := getControllerPort()
	if controllerPort == "" {
		controllerPort = ControllerPort
	}
	controllerHost := getControllerHost()
	if controllerHost == "" {
		controllerHost = ControllerHost
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", controllerHost, controllerPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	client := pb.NewCoreFileServiceClient(conn)
	return client, conn
}
