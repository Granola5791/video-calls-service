package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

func InitMpegDash(meetingID string, userID int) (*exec.Cmd, io.WriteCloser, error) {
	cmd := exec.Command(
		"ffmpeg",
		"-f", "webm",
		"-i", "pipe:0",

		// "-filter_complex",
		// "[0:v]split=3[v1][v2][v3];"+
		// 	"[v1]scale=854:480[v1out];"+
		// 	"[v2]scale=640:360[v2out];"+
		// 	"[v3]scale=426:240[v3out]",

		// "-map", "[v1out]",
		// "-map", "[v2out]",
		// "-map", "[v3out]",

		"-s", "854x480",

		"-map", "0:v",

		"-map", "0:a",
		"-c:a", "libopus",
		"-b:a", "96k",

		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-tune", "zerolatency",

		"-r", "8",

		"-g", "4",
		"-keyint_min", "4",
		"-sc_threshold", "0",

		"-profile:v", "main",

		// "-b:v:0", "1200k",
		// "-b:v:1", "750k",
		// "-b:v:2", "300k",

		"-b:v", "600k",

		"-f", "dash",
		"-ldash", "1",
		"-streaming", "1",
		"-seg_duration", "0.25",
		// "-min_seg_duration", "1000000",
		// "-frag_duration", "0.1",
		// "-frag_type", "duration",
		"-use_template", "1",
		"-use_timeline", "1",
		"-utc_timing_url", "https://time.akamai.com/?iso",
		"-window_size", "5",
		"-extra_window_size", "5",
		"-write_prft", "1",
		// "-adaptation_sets", "id=0,streams=0,1,2 id=1,streams=3",
		"-adaptation_sets", "id=0,streams=0 id=1,streams=1",
		"-remove_at_exit", "0",
		fmt.Sprintf("%s/%s/%d/stream.mpd", GetStringFromConfig("meeting.dir_path"), meetingID, userID),
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	LogFfmpeg(cmd)

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
