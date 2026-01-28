import { useEffect, useRef } from 'react';
import shaka from 'shaka-player/dist/shaka-player.ui';
import 'shaka-player/dist/controls.css';
import { StyledVideo } from '../styled-components/StyledVideos';
import { StreamConfig } from '../constants/general-contants';
import { Sleep } from '../utils/sleep';

interface DashPlayerProps {
    url: string;
}

const DashPlayer = ({ url }: DashPlayerProps) => {
    const videoRef = useRef<HTMLVideoElement>(null);
    const playerRef = useRef<shaka.Player | null>(null);

    useEffect(() => {
        try {
            startStream();
        } catch (error) {
            console.error('Error starting stream:', error);
        }

        return () => {
            stopStream();
        };
    }, []);

    const WaitUntilAvailable = async (url: string): Promise<void> => {
        let res, available = false;
        while (!available) {
            try {
                res = await fetch(url, { method: 'HEAD' });
                if (res.ok) {
                    available = true;
                }
            } catch (error) {
                console.error(error);
            }
            await Sleep(StreamConfig.checkStreamAvailabilityIntervalMs);
        }
    };

    const startStream = async () => {
        await WaitUntilAvailable(url);

        if (!videoRef.current) return;

        const player = new shaka.Player(videoRef.current);
        playerRef.current = player;

        player.configure({
            streaming: {
                lowLatencyMode: true,
                bufferingGoal: 2,       // seconds to buffer normally
                rebufferingGoal: 1,     // seconds to buffer after a stall
            },
        });


        await player.load(url);

        videoRef.current.play();
    };

    const stopStream = () => {
        if (videoRef.current?.srcObject) {
            const stream = videoRef.current.srcObject as MediaStream;
            stream.getTracks().forEach(track => track.stop());
            videoRef.current.srcObject = null;
        }

        if (playerRef.current) {
            playerRef.current.destroy();
            playerRef.current = null;
        }

        if (videoRef.current) {
            videoRef.current.pause();
            videoRef.current.src = '';
        }
    };

    return (
        <StyledVideo ref={videoRef} />
    );
};

export default DashPlayer;