import { useEffect, useRef } from 'react'
import { StreamConfig } from '../constants/general-contants';
import { StyledVideo } from '../styled-components/StyledVideos';

interface WebSocketWebCamProps {
    wsUrl: string
}

const WebSocketWebCam = ({ wsUrl}: WebSocketWebCamProps) => {
    const videoRef = useRef<HTMLVideoElement>(null);
    const wsRef = useRef<WebSocket | null>(null);
    const recorderRef = useRef<MediaRecorder | null>(null);

    useEffect(() => {
        StartStream();

        return () => {
            StopStream(); // Cleanup on unmount
        };
    }, []);

    const StartStream = async () => {
        // Ask for camera access
        const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });

        // Show preview
        if (videoRef.current) {
            videoRef.current.srcObject = stream;
            await videoRef.current.play();
        }

        // Connect to WebSocket
        const ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        // Start MediaRecorder
        const mimeType = 'video/webm; codecs=vp8,opus';
        const recorder = new MediaRecorder(stream, { mimeType });
        recorderRef.current = recorder;

        recorder.ondataavailable = (event: BlobEvent) => {
            if (event.data.size > 0) {
                event.data.arrayBuffer().then((buffer) => {
                    if (ws.readyState === WebSocket.OPEN) {
                        ws.send(buffer);
                        console.log(StreamConfig.sentChunkMsg, new Date().toISOString());
                    }
                });
            }
        };

        ws.onmessage = (event) => {
            if (event.data === StreamConfig.serverReadyMsg) {
                recorder.start(StreamConfig.chunkIntervalMs); // Send data every second
            }
        }
    };

    const StopStream = () => {
        if (recorderRef.current) {
            recorderRef.current.stop();
        }
        if (wsRef.current) {
            wsRef.current.close();
        }
        if (videoRef.current && videoRef.current.srcObject) {
            const stream = videoRef.current.srcObject as MediaStream;
            stream.getTracks().forEach(track => track.stop());
        }
    };

    return (
        <div>
            <StyledVideo
                ref={videoRef}
                autoPlay
                playsInline
                muted
            />
        </div>
    )
}

export default WebSocketWebCam