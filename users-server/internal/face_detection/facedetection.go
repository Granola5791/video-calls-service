package face_detection

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Granola5791/video-calls-service/internal/config"
)

type FaceDetectionResponse struct {
	FramesWithFace int `json:"frames_with_face"`
	TotalFrames    int `json:"total_frames"`
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

func PassedFaceDetectionThreshold(framesWithFace, totalFrames int) bool {
	if framesWithFace == 0 || totalFrames == 0 { // to avoid division by zero
		return false
	}
	minTotalFrames := config.GetIntFromConfig("face_detection.min_total_frames")
	minFaceFramesPercentage := config.GetIntFromConfig("face_detection.min_frames_with_face_percentage")
	FaceFramePercentage := 100 * framesWithFace / totalFrames
	return totalFrames >= minTotalFrames && FaceFramePercentage >= minFaceFramesPercentage
}
