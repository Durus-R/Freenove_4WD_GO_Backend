package rpc

import (
	car "Freenove_4WD_GO_Backend/Car"
	pb "Freenove_4WD_GO_Backend/dist/proto"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

type UltrasonicServer struct {
	pb.UltrasonicServer
	ultrasonic *car.Ultrasonic
}

func (s UltrasonicServer) Probe(context.Context, *empty.Empty) (*pb.UltrasonicResult, error) {
	res, err := s.ultrasonic.Probe()
	return &pb.UltrasonicResult{Result: float32(res)}, err
}
