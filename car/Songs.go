package car

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Note struct {
	Duration float64 // Beats - see https://github.com/hybridgroup/gobot/blob/v2.1.1/drivers/gpio/buzzer_driver.go#L10
	Pitch    float64 // Hz or 0 in a Pause
}

type Song []Note

func (s Song) EstimatedDuration(bpm float64) float64 {
	res := 0.0
	for _, tone := range s {
		res += tone.Duration
	}
	return 60 / bpm * res
}

func containsLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func ParseSongFile(reader io.Reader) (Song, error) {
	if reader == nil {
		return nil, errors.New("provided reader is nil")
	}

	scanner := bufio.NewScanner(reader)
	var song Song
	for scanner.Scan() {
		line := scanner.Text()
		subfields := strings.Fields(line)

		if len(subfields) < 1 || len(subfields) > 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}

		duration, err := strconv.ParseFloat(subfields[0], 64)
		if err != nil {
			return nil, fmt.Errorf("can't parse Duration: %v", err)
		}

		pitch := 0.0
		if len(subfields) > 1 {
			pitchString := subfields[1]
			if containsLetter(pitchString) {
				var ok bool
				pitch, ok = PitchMap[pitchString]
				if !ok {
					return nil, fmt.Errorf("unknown Pitch: %v", err)
				}
			} else {
				pitch, err = strconv.ParseFloat(pitchString, 64)
				if err != nil {
					return nil, fmt.Errorf("can't parse Pitch: %v", err)
				}
			}
		}

		song = append(song, Note{Duration: duration, Pitch: pitch})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %v", err)
	}

	return song, nil
}
