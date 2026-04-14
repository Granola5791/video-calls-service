package logger

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"github.com/Granola5791/video-calls-service/internal/config"
)

var ffmpegLogger *log.Logger

func InitLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	InitFfmpegLogger()
}

func InitFfmpegLogger() {
	file, err := os.Create(config.GetStringFromConfig("logger.ffmpeg_log_file"))
	if err != nil {
		log.Fatal(err)
	}
	ffmpegLogger = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

func LogFfmpeg(cmd *exec.Cmd) {
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			ffmpegLogger.Println("[ffmpeg stderr]", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			ffmpegLogger.Println("[ffmpeg stdout]", scanner.Text())
		}
	}()
}
