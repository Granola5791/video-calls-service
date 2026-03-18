package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/google/uuid"
)

type FaceDetectionResponse struct {
	FramesWithFace int `json:"frames_with_face"`
	TotalFrames    int `json:"total_frames"`
}

func ConcatenateChunks(chunks []UserVideoChunk) (*io.PipeReader, error) {
	outPipeRead, outPipeWrite := io.Pipe()
	inPipeRead, inPipeWrite := io.Pipe()

	cmd := exec.Command(
		"ffmpeg", "-y",
		"-i", "pipe:0",
		"-c", "copy",
		"-f", "webm",
		"pipe:1",
	)
	cmd.Stdin = inPipeRead
	cmd.Stdout = outPipeWrite

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	if chunks[0].ChunkNumber != 0 {
		header, err := GetHeaderUserVideoChunks(chunks[0].UserID, chunks[0].MeetingID)
		if err != nil {
			return nil, err
		}
		_, err = inPipeWrite.Write(header)
		if err != nil {
			return nil, err
		}
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

func SendvideoToFaceDetector(url string, dataPipeRead *io.PipeReader) (framesWithFace, totalFrames int, err error) {
	req, err := http.NewRequest("POST", url, dataPipeRead)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	req.Header.Set("Content-Type", "video/webm")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	defer resp.Body.Close()

	var faceDetectionResponse FaceDetectionResponse
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}
	err = json.Unmarshal(bodyBytes, &faceDetectionResponse)
	if err != nil {
		log.Println(err)
		return 0, 0, err
	}

	return faceDetectionResponse.FramesWithFace, faceDetectionResponse.TotalFrames, nil
}

func GetHeaderOfWebm(firstChunk []byte) []byte {
	clusterStart := []byte{0x1F, 0x43, 0xB6, 0x75} // start of the first cluster
	index := bytes.Index(firstChunk, clusterStart)
	return firstChunk[:index]
}

func GetHeaderUserVideoChunks(userID uint, meetingID uuid.UUID) ([]byte, error) {
	firstChunk, err := GetLatestStartChunkFromDB(meetingID, userID)
	if err != nil {
		return nil, err
	}
	return GetHeaderOfWebm(firstChunk.Chunk), nil
}

func PassedFaceDetectionThreshold(framesWithFace, totalFrames int) bool {
	if framesWithFace == 0 || totalFrames == 0 { // to avoid division by zero
		return false
	}
	minTotalFrames := GetIntFromConfig("face_detection.min_total_frames")
	minFaceFramesPercentage := GetIntFromConfig("face_detection.min_frames_with_face_percentage")
	FaceFramePercentage := 100 * framesWithFace / totalFrames
	return totalFrames >= minTotalFrames && FaceFramePercentage >= minFaceFramesPercentage
}
