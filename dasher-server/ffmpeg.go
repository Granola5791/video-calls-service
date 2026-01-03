package main

import (
	"io"
	"log"
	"os/exec"
)

func InitMpegDash(meetingID string) (*exec.Cmd, io.WriteCloser, error) {
	cmd := exec.Command(
		"ffmpeg",
		"-f", "webm",
		"-i", "pipe:0",

		"-filter_complex",
		"[0:v]split=3[v1][v2][v3];"+
			"[v1]scale=854:480[v1out];"+
			"[v2]scale=640:360[v2out];"+
			"[v3]scale=426:240[v3out]",

		"-map", "[v1out]",
		"-map", "[v2out]",
		"-map", "[v3out]",

		"-map", "0:a",
		"-c:a", "aac",
		"-b:a", "128k",

		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-tune", "zerolatency",

		"-g", "15",
		"-keyint_min", "15",
		"-sc_threshold", "0",

		"-profile:v", "main",

		"-b:v:0", "3000k",
		"-b:v:1", "1500k",
		"-b:v:2", "800k",

		"-f", "dash",
		"-use_wallclock_as_timestamps", "1",
		"-ldash", "1",
		"-streaming", "1",
		"-frag_type", "duration",
		"-frag_duration", "0.1",
		"-seg_duration", "0.5",
		"-use_template", "1",
		"-use_timeline", "1",
		"-window_size", "5",
		"-extra_window_size", "5",
		"-adaptation_sets", "id=0,streams=0,1,2 id=1,streams=3",
		"-remove_at_exit", "0",
		"meetings/" + meetingID + "/stream.mpd",
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
