import React, { useEffect, useRef } from 'react'
import { WebCamConfig } from '../constants/general-contants'

interface WebCamChunks {
    width?: string
    height?: string
}

const WebCamChunks = ({ width = WebCamConfig.defaultWidth, height = WebCamConfig.defaultHeight }: WebCamChunks) => {

    const [chunks, setChunks] = React.useState<Blob[]>([]);
    const videoRef = useRef<HTMLVideoElement>(null)
    const mediaRecorderRef = useRef<MediaRecorder | null>(null);



    const StartRecording = async () => {
        try {
            let currentStream: MediaStream | null = null;
            if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
                currentStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: false })
                if (videoRef.current) {
                    videoRef.current.srcObject = currentStream;
                }
                const mediaRecorder = new MediaRecorder(currentStream);
                mediaRecorderRef.current = mediaRecorder;

                mediaRecorder.ondataavailable = (event) => {
                    if (event.data.size > 0) {
                        setChunks((prevChunks) => [...prevChunks, event.data]);
                        console.log('Got chunk:', event.data.size, 'bytes');
                    }
                }

                mediaRecorder.start(10000);
            }
        } catch (error) {
            console.log(error);
        }

    }

    const stopRecording = () => {
        if (mediaRecorderRef.current) {
            mediaRecorderRef.current.stop();
        }
        if (videoRef.current && videoRef.current.srcObject) {
            (videoRef.current.srcObject as MediaStream).getTracks().forEach(track => track.stop());
        }
    };

    const downloadVideo = () => {
        // Combine all chunks into one Blob
        // const blob = new Blob(chunks, { type: 'video/webm' });

        // Create download link
        // const url = URL.createObjectURL(blob);
        // const a = document.createElement('a');
        // a.href = url;
        // a.download = 'recording.webm';
        // a.click();

        console.log("chunks length", chunks.length);
        for (let i = 0; i < chunks.length; i++) {
            let chunk = chunks[i];
            let blob = new Blob([chunk], { type: 'video/webm' });
            let url = URL.createObjectURL(blob);
            let a = document.createElement('a');
            a.href = url;
            a.download = `recording-${i}.webm`;
            a.click();
            URL.revokeObjectURL(url);
        }

        // Clean up
        // URL.revokeObjectURL(url);
    };


    return (
        <>
            <video ref={videoRef} autoPlay width={width} height={height} />
            <button onClick={downloadVideo}>Download Video</button>
            <button onClick={StartRecording}>Start Recording</button>
            <button onClick={stopRecording}>Stop Recording</button>
        </>
    )
}

export default WebCamChunks