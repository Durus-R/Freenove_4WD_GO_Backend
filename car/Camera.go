package car

import (
	"github.com/tractus/piCamera"
)

type Camera struct {
	cam *piCamera.PiCamera
}

func NewCamera() (*Camera, error) {
	cam, err := piCamera.New(nil, nil)
	if err != nil {
		return nil, err
	}
	err = cam.Start()
	defer cam.Stop()
	return &Camera{cam: cam}, err
}

func (c *Camera) GetFrame() ([]byte, error) {
	return c.cam.GetFrame()
}
