import { useEffect, useRef, useState } from 'react';
import shaka from 'shaka-player/dist/shaka-player.ui';
import 'shaka-player/dist/controls.css';
import { StyledVideo } from '../styled-components/StyledVideos';
import { StreamConfig } from '../constants/general-contants';
import { Sleep } from '../utils/sleep';
import { Menu, MenuItem } from '@mui/material';

interface DashPlayerProps {
    userID: string;
    url: string;
    menuOptions: { label: string; onClick: (userID: string) => void }[];
}

const DashPlayer = ({ userID, url, menuOptions }: DashPlayerProps) => {
    const videoRef = useRef<HTMLVideoElement>(null);
    const playerRef = useRef<shaka.Player | null>(null);
    const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
    const openMenu = Boolean(menuAnchorEl);

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
                res = await fetch(url, { method: 'HEAD', credentials: 'include' });
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
                bufferingGoal: 1.5,       // seconds to buffer normally
                rebufferingGoal: 0.5,     // seconds to buffer after a stall
                updateIntervalSeconds: 1,
                liveSync: {
                    enabled: true,
                    targetLatency: 0.5,
                    targetLatencyTolerance: 0.5,
                    maxPlaybackRate: 1.25,
                },
            },
            manifest: {
                updatePeriod: 0.5,
            },
        });

        console.log(player.getConfiguration());

        const networkingEngine = player.getNetworkingEngine();
        networkingEngine?.registerRequestFilter((type, request) => {
            request.allowCrossSiteCredentials = true;
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

    const handleContextMenu = (event: React.MouseEvent<HTMLVideoElement>) => {
        event.preventDefault();
        setMenuAnchorEl(event.currentTarget);
    };

    return (
        <>
            <StyledVideo
                ref={videoRef}
                onContextMenu={handleContextMenu}
            />
            <Menu anchorEl={menuAnchorEl} open={openMenu} onClose={() => setMenuAnchorEl(null)}>
                {menuOptions.map((option, index) => (
                    <MenuItem key={index} onClick={() => {option.onClick(userID)}}>
                        {option.label}
                    </MenuItem>
                ))}
            </Menu>
        </>
    );
};

export default DashPlayer;