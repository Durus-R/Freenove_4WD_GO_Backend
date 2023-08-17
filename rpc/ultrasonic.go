package rpc

import (
	car "Freenove_4WD_GO_Backend/Car"
	pb "Freenove_4WD_GO_Backend/dist/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"time"
)

type UltrasonicServer struct {
	pb.UltrasonicServer
	Ultrasonic *car.Ultrasonic
}

func (s UltrasonicServer) Probe(_ *empty.Empty, srv pb.Ultrasonic_ProbeServer) error {
	for {
		measure, err := s.Ultrasonic.Probe()
		if err != nil {
			log.Println("Failed to get Ultrasonic: ", err)
			continue
		}

		err = srv.Send(&pb.UltrasonicResult{Result: float32(measure)})
		if err != nil {
			return err
		}
		time.Sleep(time.Second / 10)
	}
}
