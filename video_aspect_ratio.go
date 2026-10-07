package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	type dimensions struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	dims := &dimensions{}

	err = json.Unmarshal(outputBuffer.Bytes(), dims)
	if err != nil {
		return "", err
	}

	if len(dims.Streams) == 1 {
		return "", errors.New("ffprobe returned no streams")
	}

	width := dims.Streams[0].Width
	height := dims.Streams[0].Height

	if height == 0 {
		return "", errors.New("invalid height dimension")
	}

	aspectRatio := float64(width) / float64(height)

	const (
		target16by9 = 16.0 / 9.0
		target9by16 = 9.0 / 16.0
		tolerance   = 0.1
	)

	if math.Abs(aspectRatio-target16by9) <= tolerance {
		return "16:9", nil
	}
	if math.Abs(aspectRatio-target9by16) <= tolerance {
		return "9:16", nil
	}

	return "other", nil
}
