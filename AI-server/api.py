import os

from fastapi import FastAPI, Request, HTTPException, BackgroundTasks
import mediapipe as mp
from mediapipe.tasks import python
from mediapipe.tasks.python import vision
import cv2
import uvicorn
from face_detector import FaceDetector
import shutil
from pathlib import Path
import tempfile
import ssl
import os
from dotenv import load_dotenv

load_dotenv()

app = FastAPI()
ssl_context = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
ssl_context.load_cert_chain(os.getenv("TLS_CERT_PATH"), os.getenv("TLS_KEY_PATH"))


@app.post("/face-detection")
async def face_detection(request: Request, BackgroundTasks: BackgroundTasks):
    frames_with_face = 0
    total_frames = 0

    video_bytes = await request.body()

    face_detector = FaceDetector()
    try:
        with tempfile.NamedTemporaryFile(suffix=".webm", delete=False) as tfile:
            tfile.write(video_bytes)
            tfile.flush()
            cap = cv2.VideoCapture(tfile.name)
            fps = cap.get(cv2.CAP_PROP_FPS) or 30

            while cap.isOpened():
                success, frame = cap.read()
                if not success:
                    break
                mp_image = mp.Image(
                    image_format=mp.ImageFormat.SRGB,
                    data=cv2.cvtColor(frame, cv2.COLOR_BGR2RGB),
                )
                timestamp_ms = int((total_frames / fps) * 1000)
                detection_result = face_detector.detect(mp_image, timestamp_ms)
                if detection_result.detections:
                    frames_with_face += 1
                total_frames += 1
            cap.release()
    finally:
        if os.path.exists(tfile.name):
            os.remove(tfile.name)
        BackgroundTasks.add_task(face_detector.restart)
    print("total_frames:", total_frames)
    print("frames_with_face:", frames_with_face)
    return {
        "frames_with_face": frames_with_face,
        "total_frames": total_frames,
    }


if __name__ == "__main__":
    uvicorn.run(
        os.getenv("APP_NAME"),
        host=os.getenv("HOST"),
        port=int(os.getenv("PORT")),
        ssl_certfile=os.getenv("TLS_CERT_PATH"),
        ssl_keyfile=os.getenv("TLS_KEY_PATH"),
    )
