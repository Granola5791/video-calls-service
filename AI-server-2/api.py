from fastapi import FastAPI, Request, HTTPException
import mediapipe as mp
from mediapipe.tasks import python
from mediapipe.tasks.python import vision
import cv2
from face_detector import FaceDetector
import shutil
from pathlib import Path

face_detector = FaceDetector()

app = FastAPI()

# Create a directory to store uploaded videos
UPLOAD_DIR = Path("uploaded_videos")
UPLOAD_DIR.mkdir(exist_ok=True)


@app.get("/hello")
def hello():
    print("hello")
    return {"message": "Hello World"}


@app.get("/face-detection")
def face_detection():
    pass


# @app.post("/upload-video/")
# async def upload_video(file:):
#     # Create a destination path
#     file_path = UPLOAD_DIR / file.filename

#     # Use shutil to save the file efficiently
#     with file_path.open("wb") as buffer:
#         shutil.copyfileobj(file.file, buffer)

#     return {
#         "filename": file.filename,
#         "content_type": file.content_type,
#         "saved_path": str(file_path),
#     }

@app.get("/upload-video")
async def upload_video(request: Request):
    """
    Receives a streamed video upload and saves it to 'file.webm'.
    This handles the data in chunks to keep memory usage low.
    """
    try:
        # Open a local file in 'write binary' mode
        with open("file.webm", "wb") as buffer:
            # Iterate through the body of the request in chunks
            # default chunk_size is 64kb
            async for chunk in request.stream():
                buffer.write(chunk)
                
        return {"status": "success", "filename": "file.webm"}
        
    except Exception as e:
        # If something goes wrong (e.g., disk full, connection lost)
        raise HTTPException(status_code=500, detail=f"Failed to save video: {str(e)}")
