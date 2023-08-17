package rpc

import (
	car "Freenove_4WD_GO_Backend/Car"
	pb "Freenove_4WD_GO_Backend/dist/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"time"
)

type ADCServer struct {
	pb.UnimplementedADCServer
	ADC *car.ADC
}

func (a *ADCServer) Stream(_ *empty.Empty, srv pb.ADC_StreamServer) error {
	for {
		battery, err := a.ADC.Battery()
		if err != nil {
			log.Println("Failed to get Battery: ", err)
			continue
		}

		idrProbe, err := a.ADC.IDR()
		if err != nil {
			log.Println("Failed to get IDR: ", err)
			continue
		}

		err = srv.Send(&pb.ADCState{
			Loaded: float32(battery),
			Left:   float32(idrProbe[0]),
			Right:  float32(idrProbe[1]),
		})

		time.Sleep(time.Second / 5)
	}

}
