package rpc

import (
	car "Freenove_4WD_GO_Backend/Car"
	pb "Freenove_4WD_GO_Backend/dist/proto"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"strings"
)

type BuzzerServer struct {
	pb.UnimplementedBuzzerServer
	Buzz         *car.Buzzer
	currentSong  car.Song
	signalFinish chan struct{}
	isClosing    bool
}

func (s *BuzzerServer) On(context.Context, *empty.Empty) (*empty.Empty, error) {
	err := s.Buzz.On()
	return nil, err
}

func (s *BuzzerServer) Off(context.Context, *empty.Empty) (*empty.Empty, error) {
	err := s.Buzz.Off()
	return nil, err
}

func (s *BuzzerServer) Toggle(context.Context, *empty.Empty) (*empty.Empty, error) {
	err := s.Buzz.Toggle()
	return nil, err
}

func (s *BuzzerServer) SetBPM(_ context.Context, bpm *pb.SetBPMRequest) (*empty.Empty, error) {
	err := s.Buzz.SetBPM(float64(bpm.GetBpm()))
	return nil, err
}

func (s *BuzzerServer) GetBPM(context.Context, *empty.Empty) (*pb.GetBPMRequest, error) {
	return &pb.GetBPMRequest{Bpm: float32(s.Buzz.GetBPM())}, nil
}

func (s *BuzzerServer) CalculateDuration(_ context.Context, sng *pb.Song) (*pb.SongDuration, error) {
	song := importSong(sng)
	dur := song.EstimatedDuration(s.Buzz.GetBPM())
	return &pb.SongDuration{Length: float32(dur)}, nil

}

func (s *BuzzerServer) ParseSong(_ context.Context, stringPayload *pb.SongStringPayload) (*pb.Song, error) {
	reader := strings.NewReader(stringPayload.GetPayload())
	parsed, err := car.ParseSongFile(reader)
	if err != nil {
		return &pb.Song{}, err
	}

	return exportSong(parsed), nil
}

func importSong(sng *pb.Song) car.Song {
	song := car.Song{}
	for _, n := range sng.GetNotes() {
		song = append(song, car.Note{
			Duration: float64(n.Duration),
			Pitch:    float64(n.Pitch),
		})
	}
	return song
}

func exportSong(song car.Song) *pb.Song {
	var notes []pb.Note
	for _, n := range song {
		notes = append(notes, pb.Note{
			Duration: float32(n.Duration),
			Pitch:    float32(n.Pitch),
		})
	}
	notesPointers := make([]*pb.Note, len(notes))
	for i := range notes {
		notesPointers[i] = &notes[i]
	}
	return &pb.Song{Notes: notesPointers}
}

func (s *BuzzerServer) AsyncPlaySong(_ context.Context, sng *pb.Song) (*empty.Empty, error) {
	s.currentSong = importSong(sng)
	s.signalFinish = make(chan struct{})
	go func() {
		err := s.Buzz.PlaySong(s.currentSong, s.signalFinish)
		if err != nil {
			log.Println("Error at playing song: ", err)
		}
		s.signalFinish = nil
		s.currentSong = car.Song{}
		s.isClosing = false
	}()
	return nil, nil
}

func (s *BuzzerServer) DoesSongStillPlay(context.Context, *empty.Empty) (*pb.SongStatus, error) {
	return &pb.SongStatus{
		IsPlaying: s.currentSong != nil,
	}, nil
}

func (s *BuzzerServer) StopSong(context.Context, *empty.Empty) (*empty.Empty, error) {
	if s.signalFinish != nil {
		if s.isClosing == false {
			s.isClosing = true
			close(s.signalFinish)
		}
	}
	return nil, nil
}

func (s *BuzzerServer) GetSong(context.Context, *empty.Empty) (*pb.Song, error) {
	return exportSong(s.currentSong), nil
}
