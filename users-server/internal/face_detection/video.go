package face_detection

import (
	"bytes"
	"io"
	"log"
	"os/exec"

	"github.com/Granola5791/video-calls-service/internal/db"
	"github.com/google/uuid"
)

func ConcatenateChunks(chunks []db.UserVideoChunk) (*io.PipeReader, error) {
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

func GetHeaderOfWebm(firstChunk []byte) []byte {
	clusterStart := []byte{0x1F, 0x43, 0xB6, 0x75} // start of the first cluster
	index := bytes.Index(firstChunk, clusterStart)
	return firstChunk[:index]
}

func GetHeaderUserVideoChunks(userID uint, meetingID uuid.UUID) ([]byte, error) {
	firstChunk, err := db.GetLatestStartChunk(meetingID, userID)
	if err != nil {
		return nil, err
	}
	return GetHeaderOfWebm(firstChunk.Chunk), nil
}
