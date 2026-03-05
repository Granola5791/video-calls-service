import mediapipe as mp
from mediapipe.tasks import python
from mediapipe.tasks.python import vision
import cv2
import threading


class FaceDetector:
    def __init__(self):
        self.base_options = python.BaseOptions(model_asset_path="detector.tflite")
        self.options = vision.FaceDetectorOptions(
            base_options=self.base_options, running_mode=vision.RunningMode.VIDEO
        )
        self.detector = vision.FaceDetector.create_from_options(self.options)
        self.lock = threading.Lock()

    def restart(self):
        with self.lock:
            self.detector.close()
            self.detector = vision.FaceDetector.create_from_options(self.options)

    def detect(self, mp_image, timestamp_ms):
        with self.lock:
            return self.detector.detect_for_video(mp_image, timestamp_ms)