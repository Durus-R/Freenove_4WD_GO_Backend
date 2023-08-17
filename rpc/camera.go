package rpc

import (
	car "Freenove_4WD_GO_Backend/car"
	pb "Freenove_4WD_GO_Backend/dist/proto"
	"log"
	"time"
)

type CameraServer struct {
	pb.UnimplementedCameraServer
	Camera *car.Camera
}

func (c *CameraServer) StreamCamera(fpsRequest *pb.FramesPerSecond, srv pb.Camera_StreamCameraServer) error {
	fps := time.Duration(fpsRequest.GetFPS())
	for {
		frame, err := c.Camera.GetFrame()
		if err != nil {
			log.Println("Failed to read frame: ", err)
			continue
		}

		err = srv.Send(&pb.CameraStream{Frame: frame})
		if err != nil {
			return err
		}

		time.Sleep(time.Second / fps) // Delay between frames
	}
}
