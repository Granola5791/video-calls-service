package main

import (
	"io"
	"log"
	"os/exec"
)

func InitMpegDash() (*exec.Cmd, io.WriteCloser, error) {
	cmd := exec.Command(
		"ffmpeg",
		"-f", "webm",
		"-i", "pipe:0",

		"-filter_complex",
		"[0:v]split=3[v1][v2][v3];"+
			"[v1]scale=1920:1080[v1out];"+
			"[v2]scale=1280:720[v2out];"+
			"[v3]scale=854:480[v3out]",

		"-map", "[v1out]",
		"-map", "[v2out]",
		"-map", "[v3out]",

		"-c:v", "libx264",
		"-preset", "veryfast",
		"-tune", "zerolatency",

		"-g", "30",
		"-keyint_min", "30",
		"-sc_threshold", "0",

		"-profile:v", "main",

		"-b:v:0", "5000k",
		"-b:v:1", "3000k",
		"-b:v:2", "1500k",

		"-f", "dash",
		"-seg_duration", "1",
		"-use_template", "1",
		"-use_timeline", "1",
		"-window_size", "5",
		"-extra_window_size", "5",
		"-adaptation_sets", "id=0,streams=0 id=1,streams=1 id=2,streams=2",
		"data/stream.mpd",
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	return cmd, stdin, nil
}

func PrepareForMpegDash(stdin io.WriteCloser, video []byte) {
	_, err := stdin.Write(video)
	if err != nil {
		log.Println(err)
	}
}
