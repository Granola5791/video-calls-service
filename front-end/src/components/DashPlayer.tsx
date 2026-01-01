import React, { useRef, useState } from 'react';
import shaka from 'shaka-player/dist/shaka-player.ui';
import 'shaka-player/dist/controls.css';
import { DasherServerAddressHttp, ApiEndpoints } from '../constants/backend-constants';

const ShakaPlayer = () => {
    const videoRef = useRef<HTMLVideoElement>(null);
    const containerRef = useRef<HTMLDivElement>(null);
    const playerRef = useRef<shaka.Player | null>(null);
    const [isPlaying, setIsPlaying] = useState(false);

    const startStream = async () => {
        if (!videoRef.current || !containerRef.current) return;

        const player = new shaka.Player(videoRef.current);
        playerRef.current = player;

        player.configure({
            streaming: {
                lowLatencyMode: true,
                bufferingGoal: 1,
                rebufferingGoal: 0.5,
                startupBufferingGoal: 0.1,
                stallEnabled: false,
                liveSync: {
                    enabled: true,
                    targetLatency: 1.0,
                    maxLatency: 2.0,
                },
            },
        });

        await player.load(DasherServerAddressHttp + ApiEndpoints.getStream);

        videoRef.current.play();
        setIsPlaying(true);
    };

    const stopStream = () => {
        if (playerRef.current) {
            playerRef.current.destroy();
            playerRef.current = null;
            setIsPlaying(false);
        }
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

                <div ref={containerRef} className="w-full">
                    <video ref={videoRef} className="w-full" />
                </div>
            </div>
        </div>
    );
};

export default ShakaPlayer;