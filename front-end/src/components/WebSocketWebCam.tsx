import { useEffect, useRef } from 'react'
import { CenteredColumn } from '../styled-components/StyledBoxes';
import { StreamConfig } from '../constants/general-contants';

interface WebSocketWebCamProps {
    wsUrl: string
    width?: string
    height?: string
}

const WebSocketWebCam = ({ wsUrl, width, height }: WebSocketWebCamProps) => {
    const videoRef = useRef<HTMLVideoElement>(null);
    const wsRef = useRef<WebSocket | null>(null);
    const recorderRef = useRef<MediaRecorder | null>(null);

    useEffect(() => {
        return () => {
            StopStream(); // Cleanup on unmount
        };
    }, []);

    const StartStream = async () => {
        // Connect to WebSocket
        const ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        // Ask for camera access
        const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: false });

        // Show preview
        if (videoRef.current) {
            videoRef.current.srcObject = stream;
            await videoRef.current.play();
        }

        // Start MediaRecorder
        const mimeType = 'video/webm; codecs=vp8';
        const recorder = new MediaRecorder(stream, { mimeType });
        recorderRef.current = recorder;

        recorder.ondataavailable = (event: BlobEvent) => {
            if (event.data.size > 0) {
                event.data.arrayBuffer().then((buffer) => {
                    if (ws.readyState === WebSocket.OPEN) {
                        ws.send(buffer);
                        console.log('Sent video chunk of size:', buffer.byteLength, "at", new Date().toISOString());
                    }
                });
            }
        };

        ws.onmessage = (event) => {
            if (event.data === StreamConfig.serverReadyMsg) {
                recorder.start(StreamConfig.chunkIntervalMs); // Send data every second
            }
        }
    }

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
    }

    return (
        <CenteredColumn>
            <video
                ref={videoRef}
                autoPlay
                width={width}
                height={height}
            />
            <button onClick={StopStream}>Stop</button>
            <button onClick={StartStream}>Start</button>
        </CenteredColumn>
    )
}

export default WebSocketWebCam