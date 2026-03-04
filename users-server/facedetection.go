package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleFaceDetection(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))
	userID := c.GetInt(GetStringFromConfig("jwt.user_id_name"))
	videoChunks, err := GetUserVideoChunksFromDB(meetingID, uint(userID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	go MarkUserVideoChunksAsVisitedInDB(
		meetingID,
		uint(userID),
		videoChunks[0].ChunkNumber,
		videoChunks[len(videoChunks)-1].ChunkNumber,
	)

	outputPipeRead, err := ConcatenateChunks(videoChunks)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = SendvideoToFaceDetector(outputPipeRead)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
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

func SendvideoToFaceDetector(dataPipeRead *io.PipeReader) error {
	req, err := http.NewRequest("GET", "http://127.0.0.1:8000//upload-video", dataPipeRead)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "video/webm")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}