import React, { useRef } from 'react'
import { WebCamConfig } from '../constants/general-contants'

interface WebCamProps {
    width?: string,
    height?: string,
}

const WebCam = ({ width = WebCamConfig.defaultWidth, height = WebCamConfig.defaultHeight }: WebCamProps) => {

    const videoRef = useRef<HTMLVideoElement>(null)

    React.useEffect(() => {
        let currentStream: MediaStream | null = null;
        if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
            navigator.mediaDevices.getUserMedia({ video: true, audio: false })
                .then((newStream) => {
                    currentStream = newStream;

                    if (videoRef.current) {
                        videoRef.current.srcObject = currentStream;
                    }
                })
                .catch((error) => {
                    console.log(error);
                });
        }

        return () => {
            const streamToStop = currentStream;
            if (streamToStop) {
                streamToStop.getTracks().forEach((track) => {
                    track.stop();
                });
            }
        };
    }, []);


    return (
        <>
            <video ref={videoRef} autoPlay width={width} height={height} />
        </>
    )
}

export default WebCam