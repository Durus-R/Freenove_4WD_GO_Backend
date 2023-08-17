package rpc

import (
	pb "Freenove_4WD_GO_Backend/dist/proto"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"net/http"
	"time"
)

type CarServer struct {
	pb.UnimplementedCarServer
}

func (s *CarServer) Ping(context.Context, *empty.Empty) (*pb.Pong, error) {
	result, err := googlePing()
	return &pb.Pong{
		PingGoogle: uint32(result),
	}, err
}

func googlePing() (time.Duration, error) {
	timeout := 30 * time.Second
	client := http.Client{
		Timeout: timeout,
	}
	before := time.Now()
	_, err := client.Head("https://www.google.com")

	result := time.Since(before) * time.Millisecond

	return result, err
}
