import React, { useEffect, useRef, useState } from 'react';
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
    const containerRef = useRef<HTMLDivElement>(null);
    const playerRef = useRef<shaka.Player | null>(null);
    const [isPlaying, setIsPlaying] = useState(false);

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
        
        if (!videoRef.current || !containerRef.current) return;

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
        setIsPlaying(true);
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

        setIsPlaying(false);
    };

    return (
        <div className="p-8 bg-gray-900 min-h-screen">
            <div className="w-full max-w-4xl mx-auto">
                {!isPlaying && (
                    <button
                        onClick={startStream}
                        className="mb-4 px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition"
                    >
                        Start Stream
                    </button>
                )}

                {isPlaying && (
                    <button
                        onClick={stopStream}
                        className="mb-4 px-6 py-3 bg-red-600 hover:bg-red-700 text-white font-semibold rounded-lg transition"
                    >
                        Stop Stream
                    </button>
                )}

                <div ref={containerRef}>
                    <StyledVideo ref={videoRef} />
                </div>
            </div>
        </div>
    );
};

export default DashPlayer;