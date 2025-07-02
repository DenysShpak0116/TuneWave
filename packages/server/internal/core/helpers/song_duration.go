package helpers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/faiface/beep/mp3"
)

type ReadSeekCloser struct {
	*bytes.Reader
}

func (r *ReadSeekCloser) Close() error {
	return nil
}

func GetAudioDuration(file multipart.File) (time.Duration, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %w", err)
	}

	r := &ReadSeekCloser{Reader: bytes.NewReader(data)}

	streamer, format, err := mp3.Decode(r)
	if err != nil {
		return 0, fmt.Errorf("failed to decode mp3: %w", err)
	}
	defer streamer.Close()

	samples := streamer.Len()
	if samples <= 0 {
		return 0, fmt.Errorf("invalid number of samples: %d", samples)
	}

	duration := time.Duration(samples) * time.Second / time.Duration(format.SampleRate)
	return duration, nil
}
