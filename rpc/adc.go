package rpc

import (
	car "Freenove_4WD_GO_Backend/Car"
	pb "Freenove_4WD_GO_Backend/dist/proto"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

type ADCServer struct {
	pb.UnimplementedADCServer
	ADC *car.ADC
}

func (s *ADCServer) Battery(_ context.Context, _ *empty.Empty) (*pb.BatteryState, error) {
	batteryInPercent, err := s.ADC.Battery()

	return &pb.BatteryState{Loaded: float32(batteryInPercent)}, err
}

func (s *ADCServer) IDR(_ context.Context, _ *empty.Empty) (*pb.IDRState, error) {
	idrMeasure, err := s.ADC.IDR()

	return &pb.IDRState{
		Left:  float32(idrMeasure[0]),
		Right: float32(idrMeasure[1]),
	}, err
}
