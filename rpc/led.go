package rpc

import (
	car "Freenove_4WD_GO_Backend/Car"
	pb "Freenove_4WD_GO_Backend/dist/proto"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"sync"
	"time"
)

type LEDServer struct {
	pb.UnimplementedLEDServer
	led          *car.RGBStrip
	signalFinish chan struct{}
	closeLocker  sync.Mutex
}

func (s *LEDServer) StopEffect(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	var err error
	s.closeLocker.Lock()
	if s.signalFinish != nil {
		close(s.signalFinish)
		time.Sleep(200 * time.Millisecond)
		err = s.led.Black()
	}
	s.closeLocker.Unlock()
	return nil, err
}

func makeColor(col *pb.Color) uint32 {
	return car.RgbToColor(int(col.GetRed()), int(col.GetGreen()), int(col.GetBlue()))
}

func (s *LEDServer) StartColorWipe(_ context.Context, col *pb.Color) (*empty.Empty, error) {
	s.signalFinish = make(chan struct{})
	color := makeColor(col)
	err := s.led.ColorWipe(color)
	s.signalFinish = nil
	return nil, err
}

func (s *LEDServer) StartTheaterChase(_ context.Context, col *pb.Color) (*empty.Empty, error) {
	s.signalFinish = make(chan struct{})
	color := makeColor(col)
	go func() {
		_ = s.led.TheaterChase(color, s.signalFinish)
		s.signalFinish = nil
	}()
	return nil, nil
}

func (s *LEDServer) StartRainbow(_ context.Context, _ *pb.Color) (*empty.Empty, error) {
	s.signalFinish = make(chan struct{})
	go func() {
		_ = s.led.Rainbow(s.signalFinish)
		s.signalFinish = nil
	}()
	return nil, nil
}

func (s *LEDServer) StartRainbowCycle(_ context.Context, _ *pb.Color) (*empty.Empty, error) {
	s.signalFinish = make(chan struct{})
	go func() {
		_ = s.led.RainbowCycle(s.signalFinish)
		s.signalFinish = nil
	}()
	return nil, nil
}

func (s *LEDServer) StartTheaterChaseRainbow(_ context.Context, _ *pb.Color) (*empty.Empty, error) {
	s.signalFinish = make(chan struct{})
	go func() {
		_ = s.led.TheaterChaseRainbow(s.signalFinish)
		s.signalFinish = nil
	}()
	return nil, nil
}

func makeColors(colors *pb.Colors) [8]uint32 {
	var res []uint32
	items := colors.GetColors()
	for _, col := range items {
		res = append(res, makeColor(col))
	}
	return [8]uint32{res[0], res[1], res[2], res[3], res[4], res[5], res[6], res[7]}
}

func (s *LEDServer) ApplyCustomColors(_ context.Context, colors *pb.Colors) (*empty.Empty, error) {
	err := s.led.ApplyColors(makeColors(colors))
	return nil, err
}

func (s *LEDServer) IsDark(_ context.Context, _ *empty.Empty) (*pb.IsDarkResult, error) {
	return &pb.IsDarkResult{Dark: s.signalFinish == nil}, nil
}

func (s *LEDServer) EffectIsRunning(_ context.Context, _ *empty.Empty) (*pb.LockResult, error) {
	return &pb.LockResult{Locked: s.led.IsLocked()}, nil
}
