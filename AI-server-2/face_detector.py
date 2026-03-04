import mediapipe as mp
from mediapipe.tasks import python
from mediapipe.tasks.python import vision
import cv2


class FaceDetector:
    def __init__(self):
        self.base_options = python.BaseOptions(model_asset_path="detector.tflite")
        self.options = vision.FaceDetectorOptions(
            base_options=self.base_options, running_mode=vision.RunningMode.VIDEO
        )
        self.face_detector = vision.FaceDetector.create_from_options(self.options)
