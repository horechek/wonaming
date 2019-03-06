package main

import (
	"flag"
	"github.com/wothing/wonaming/consul"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wothing/wonaming/etcdv3/example/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var addr = "127.0.0.1:2379"

func main() {
	flag.StringVar(&addr, "addr", addr, "addr to lis")
	flag.Parse()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	defer s.GracefulStop()

	pb.RegisterHelloServiceServer(s, &hello{})

	go func() {
		consul.Register("hello", "127.0.0.1", 2379, "127.0.0.1:8500", 5*time.Second, 10)
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch

		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

type hello struct {
}

func (*hello) Echo(ctx context.Context, req *pb.Payload) (*pb.Payload, error) {
	req.Data = req.Data + ", from:" + addr
	return req, nil
}
