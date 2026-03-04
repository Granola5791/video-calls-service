package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
)

type FaceDetectionResponse struct {
	FramesWithFace int `json:"frames_with_face"`
}

func ConcatenateChunks(chunks []UserVideoChunk) (*io.PipeReader, error) {
	outPipeRead, outPipeWrite := io.Pipe()
	inPipeRead, inPipeWrite := io.Pipe()

	cmd := exec.Command(
		"ffmpeg",
		"-i",
		"pipe:0",
		"-c",
		"copy",
		"-f",
		"webm",
		"pipe:1",
	)
	cmd.Stdin = inPipeRead
	cmd.Stdout = outPipeWrite

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	go func() {
		defer inPipeWrite.Close()
		for _, chunk := range chunks {
			_, err := inPipeWrite.Write(chunk.Chunk)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Println(err)
		}
		outPipeWrite.CloseWithError(err)
	}()

	return outPipeRead, nil
}

func SendvideoToFaceDetector(url string, dataPipeRead *io.PipeReader) (framesWithFace int, err error) {
	req, err := http.NewRequest("GET", url, dataPipeRead)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "video/webm")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var faceDetectionResponse FaceDetectionResponse
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(bodyBytes, &faceDetectionResponse)
	if err != nil {
		return 0, err
	}

	return faceDetectionResponse.FramesWithFace, nil
}
