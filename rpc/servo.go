package rpc

import (
	car "Freenove_4WD_GO_Backend/car"
	pb "Freenove_4WD_GO_Backend/dist/proto"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServoServer struct {
	pb.UnimplementedServoServer
	MC *car.MotorController
}

func (s *ServoServer) SetVerticalAngle(_ context.Context, angle *pb.Angle) (*emptypb.Empty, error) {
	s.MC.SetAngle(0, uint16(angle.GetAngle()))
	return nil, nil
}

func (s *ServoServer) SetHorizontalAngle(_ context.Context, angle *pb.Angle) (*emptypb.Empty, error) {
	s.MC.SetAngle(1, uint16(angle.GetAngle()))
	return nil, nil
}
