package rpc

import (
	car "Freenove_4WD_GO_Backend/car"
	pb "Freenove_4WD_GO_Backend/dist/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

type ADCServer struct {
	pb.UnimplementedADCServer
	ADC *car.ADC
}

func (a *ADCServer) Stream(_ *emptypb.Empty, srv pb.ADC_StreamServer) error {
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
