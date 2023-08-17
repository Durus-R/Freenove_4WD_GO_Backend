package rpc

import (
	car "Freenove_4WD_GO_Backend/Car"
	pb "Freenove_4WD_GO_Backend/dist/proto"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MotorServer struct {
	pb.UnimplementedMotorServer
	MC                *car.MotorController
	currentMotorModel car.Direction
}

func (s *MotorServer) Forward(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	s.MC.SetDirection(car.GetDirectionForward())
	s.currentMotorModel = car.GetDirectionForward()
	return nil, nil
}

func (s *MotorServer) Backward(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	s.MC.SetDirection(car.GetDirectionBackward())
	s.currentMotorModel = car.GetDirectionBackward()
	return nil, nil
}

func (s *MotorServer) Left(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	s.MC.SetDirection(car.GetDirectionLeft())
	s.currentMotorModel = car.GetDirectionLeft()
	return nil, nil
}

func (s *MotorServer) Right(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	s.MC.SetDirection(car.GetDirectionRight())
	s.currentMotorModel = car.GetDirectionRight()
	return nil, nil
}

func (s *MotorServer) Halt(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	s.MC.SetDirection(car.GetDirectionHalt())
	s.currentMotorModel = car.GetDirectionHalt()
	return nil, nil
}

func (s *MotorServer) SetMotorModel(_ context.Context, model *pb.MotorModel) (*emptypb.Empty, error) {
	direction := car.Direction{
		LeftUp:   int(model.LeftUp),
		LeftLow:  int(model.LeftLow),
		RightUp:  int(model.RightUp),
		RightLow: int(model.RightLow),
	}
	s.MC.SetDirection(direction)
	s.currentMotorModel = direction
	return nil, nil
}

func (s *MotorServer) GetMotorModel(_ context.Context, _ *emptypb.Empty) (*pb.MotorModel, error) {
	currentMM := s.currentMotorModel
	return &pb.MotorModel{
		LeftUp:   int32(currentMM.LeftUp),
		LeftLow:  int32(currentMM.LeftLow),
		RightUp:  int32(currentMM.RightUp),
		RightLow: int32(currentMM.RightLow),
	}, nil
}
