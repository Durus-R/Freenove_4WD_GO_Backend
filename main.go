package main

import (
	car "Freenove_4WD_GO_Backend/Car"
	pb "Freenove_4WD_GO_Backend/dist/proto"
	"Freenove_4WD_GO_Backend/rpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/user"
)

// TODO: Ultrasonic, Camera
func main() {
	usr, _ := user.Current()
	if usr.Name != "root" {
		log.Fatal("Please restart this system with root access")
	}

	s := grpc.NewServer()
	pb.RegisterUltrasonicServer(s, &rpc.UltrasonicServer{
		Ultrasonic: car.NewUltrasonic(),
	})

	pb.RegisterADCServer(s, &rpc.ADCServer{
		ADC: car.NewADC(),
	})

	cam, err := car.NewCamera()
	if err != nil {
		log.Fatal("Failed to create Camera: ", err)
	}
	pb.RegisterCameraServer(s, &rpc.CameraServer{
		Camera: cam,
	})

	buzz, err := car.NewBuzzer()
	if err != nil {
		log.Fatal("Failed to create buzzer: ", err)
	}
	pb.RegisterBuzzerServer(s, &rpc.BuzzerServer{Buzz: buzz})

	pb.RegisterCarServer(s, &rpc.CarServer{})

	pb.RegisterLEDServer(s, &rpc.LEDServer{LED: car.NewRGBStrip()})

	mc, err := car.NewMotorController()
	if err != nil {
		log.Fatal("Cannot create Motor, ", err)
	}
	pb.RegisterMotorServer(s, &rpc.MotorServer{
		MC: mc,
	})

	pb.RegisterServoServer(s, &rpc.ServoServer{
		MC: mc,
	})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	err = s.Serve(lis)
	if err != nil {
		log.Fatal("Error at serving: ", err)
	}

}
