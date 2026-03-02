import { useEffect, useRef, useState } from 'react';
import shaka from 'shaka-player/dist/shaka-player.ui';
import 'shaka-player/dist/controls.css';
import { StyledVideo } from '../styled-components/StyledVideos';
import { StreamConfig } from '../constants/general-contants';
import { Sleep } from '../utils/sleep';
import { Menu, MenuItem } from '@mui/material';
import { ErrorMsgs } from '../constants/general-contants';
import { SmallMuteButton } from '../styled-components/StyledButtons';
import { MutedIcon, SmallMutedIcon, SmallUnmutedIcon, UnmutedIcon } from '../styled-components/StyledIcons';

interface DashPlayerProps {
    userID: string;
    url: string;
    menuOptions: { label: string; onClick: (userID: string) => void }[];
}

const DashPlayer = ({ userID, url, menuOptions }: DashPlayerProps) => {
    const videoRef = useRef<HTMLVideoElement>(null);
    const playerRef = useRef<shaka.Player | null>(null);
    const [menuAnchorPos, setMenuAnchorPos] = useState<{ top: number, left: number } | undefined>(undefined);
    const openMenu = Boolean(menuAnchorPos);
    const active = useRef(false);
    const [muted, setMuted] = useState(false);

    useEffect(() => {
        try {
            active.current = true;
            startStream();
        } catch (error) {
            console.error(ErrorMsgs.cantStartStream, error);
        }

        return () => {
            try {
                active.current = false;
                stopStream();
            } catch (error) {
                console.error(ErrorMsgs.cantStopStream, error);
            }
        };
    }, []);

    const WaitUntilAvailable = async (url: string): Promise<void> => {
        let res, available = false;
        while (!available && active.current) {
            try {
                res = await fetch(url, {
                    method: 'HEAD',
                    credentials: 'include',
                });
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

        if (!videoRef.current || !active.current) return;

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

        try {
            await player.load(url);
        } catch (error) {
            console.error(ErrorMsgs.cantLoadStream, error);
            startStream();
        }
        videoRef.current?.play();


    };

    const stopStream = () => {
        try {
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
        } catch (error) {
            console.error(ErrorMsgs.cantStopStream, error);
        }
    };

    const handleContextMenu = (event: React.MouseEvent<HTMLVideoElement>) => {
        if (menuOptions.length === 0) {
            return;
        }
        event.preventDefault();
        setMenuAnchorPos({ top: event.clientY, left: event.clientX });
    };

    return (
        <>
            <StyledVideo
                ref={videoRef}
                onContextMenu={handleContextMenu}
                muted={muted}
            />
            <SmallMuteButton onClick={() => setMuted(!muted)}>
                {muted ? <SmallMutedIcon /> : <SmallUnmutedIcon />}
            </SmallMuteButton>
            <Menu
                open={openMenu}
                anchorReference='anchorPosition'
                anchorPosition={menuAnchorPos}
                onClose={() => setMenuAnchorPos(undefined)}
            >
                {menuOptions.map((option, index) => (
                    <MenuItem key={index} onClick={() => { option.onClick(userID) }}>
                        {option.label}
                    </MenuItem>
                ))}
            </Menu>
        </>
    );
};

export default DashPlayer;