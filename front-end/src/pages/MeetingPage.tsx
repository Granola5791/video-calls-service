import React, { useEffect } from 'react'
import { useParams } from 'react-router-dom';
import { ApiEndpoints, BackendAddressHttp, BackendAddressWS, CallEventTypes, DasherServerAddressHttp, DasherServerAddressWS, HttpStatusCodes, SetUrlParams } from '../constants/backend-constants';
import DashPlayer from '../components/DashPlayer';
import { StyledMeetingGrid } from '../styled-components/StyledBoxes';
import { StyledMeetingGridTile, StyledVideo } from '../styled-components/StyledVideos';
import OneButtonPopUp from '../components/OneButtonPopUp';
import { useNavigation } from '../utils/navigation';
import { LocalStorage, MeetingConfig, StreamConfig } from '../constants/general-contants';
import { HostOptions, MeetingExitText } from '../constants/hebrew-constants';
import { StyledMeetingFooter } from '../styled-components/StyledFooters';


const NormalizeMeetingIDs = (meetingIDs: unknown): string[] => {
    return Array.isArray(meetingIDs) ? meetingIDs : [];
}

const MeetingPage = () => {

    const ID = localStorage.getItem(LocalStorage.id);
    const { meetingID } = useParams();
    const streamWsRef = React.useRef<WebSocket | null>(null);
    const notificationsWsRef = React.useRef<WebSocket | null>(null);
    const clientVideoRef = React.useRef<HTMLVideoElement>(null);
    const recorderRef = React.useRef<MediaRecorder | null>(null);
    const [participantsIDs, setParticipantsIDs] = React.useState<string[]>([]);
    const [meetingState, setMeetingState] = React.useState(MeetingConfig.meetingState.none);
    const keepAliveIntervalIDRef = React.useRef(0);
    const [isHost, setIsHost] = React.useState(false);
    const [hostOptions, setHostOptions] = React.useState<{ label: string, onClick: (userID: string) => void }[]>([]);
    const {
        goToHome,
    } = useNavigation();

    useEffect(() => {
        if (meetingID) {
            JoinMeeting(meetingID);
        }

        return () => {
            CloseNotificationsConnection();
            StopStream();
            if (keepAliveIntervalIDRef.current !== 0) {
                clearInterval(keepAliveIntervalIDRef.current);
            }
        };
    }, []);

    const CloseNotificationsConnection = () => {
        if (notificationsWsRef.current) {
            notificationsWsRef.current.close();
        }
    }

    const SendKeepAlive = async () => {
        const res = await fetch(BackendAddressHttp + SetUrlParams(ApiEndpoints.keepAlive, meetingID), {
            method: 'POST',
            credentials: 'include',
        })
        if (!res.ok) {
            if(res.status === HttpStatusCodes.Unauthorized) {
                setMeetingState(MeetingConfig.meetingState.kicked)
            }
            throw new Error(res.statusText);
        }
    }

    const ContinuouslySendKeepAlive = () => {
        const intervalID = setInterval(SendKeepAlive, MeetingConfig.keepAliveIntervalMs);
        keepAliveIntervalIDRef.current = intervalID;
    }

    const ActivateHostOptions = () => {
        setHostOptions([
            { label: HostOptions.kick, onClick: KickFromMeeting },
        ]);
    }

    const StartStreaming = async (wsUrl: string) => {
        // Ask for camera access
        const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });

        if (clientVideoRef.current?.srcObject) {
            throw new Error('Already streaming');
        }
        // Show preview
        if (clientVideoRef.current) {
            clientVideoRef.current.srcObject = stream;
            await clientVideoRef.current.play();
        }

        // Connect to WebSocket
        const ws = new WebSocket(wsUrl);
        streamWsRef.current = ws;

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

    const JoinMeetingBackend = async (meetingID: string) => {
        let res = await fetch(BackendAddressHttp + SetUrlParams(ApiEndpoints.joinMeeting, meetingID), {
            method: 'POST',
            credentials: 'include',
        })
        if (res.ok) {
            const data = await res.json();
            const participants = data.participants
            const is_host = data.is_host
            setParticipantsIDs(NormalizeMeetingIDs(participants));
            setIsHost(is_host);
            if (is_host) {
                ActivateHostOptions();
            }
        } else if (res.status === HttpStatusCodes.NotFound) {
            setMeetingState(MeetingConfig.meetingState.wrongID);
        }
        if (!res.ok) {
            throw new Error(res.statusText);
        }
    }

    const SubscribeToMeetingUpdates = async (meetingID: string) => {
        const ws = new WebSocket(BackendAddressWS + SetUrlParams(ApiEndpoints.getCallNotifications, meetingID));
        notificationsWsRef.current = ws;
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            const participantID = data.participant_id
            const eventType = data.event

            switch (eventType) {
                case CallEventTypes.participantJoined:
                    setParticipantsIDs((prevIDs) => prevIDs.includes(participantID) ? prevIDs : [...prevIDs, participantID]);
                    break;

                case CallEventTypes.participantLeft:
                    setParticipantsIDs((prevIDs) => prevIDs.filter((id) => id !== participantID));
                    break;
            }
        };

        ws.onclose = () => {
            notificationsWsRef.current = null;
        };
    };

    const JoinMeetingDasher = async (meetingID: string) => {
        let res = await fetch(DasherServerAddressHttp + SetUrlParams(ApiEndpoints.joinMeeting, meetingID), {
            method: 'POST',
            credentials: 'include',
        })
        if (!res.ok) {
            return;
        }
    }

    const JoinMeeting = async (meetingID: string) => {
        try {
            await JoinMeetingBackend(meetingID);
            await JoinMeetingDasher(meetingID);
            await StartStreaming(DasherServerAddressWS + SetUrlParams(ApiEndpoints.startStream, meetingID));
            await SubscribeToMeetingUpdates(meetingID);
            ContinuouslySendKeepAlive();
            setMeetingState(MeetingConfig.meetingState.active);
        } catch (error) {
            console.error(error);
        }
    }

    const StopStream = () => {
        if (recorderRef.current) {
            recorderRef.current.stop();
        }
        if (streamWsRef.current) {
            streamWsRef.current.close();
        }
        if (clientVideoRef.current && clientVideoRef.current.srcObject) {
            const stream = clientVideoRef.current.srcObject as MediaStream;
            stream.getTracks().forEach(track => track.stop());
        }
    };

    const LeaveMeeting = async () => {
        const LeaveMeetingFrontend = async (state: number) => {
            setMeetingState(state);
            StopStream();
            setTimeout(() => { goToHome() }, MeetingConfig.exitWaitTimeMs)
        }

        const LeaveMeetingDasher = () => {
            notificationsWsRef.current?.close();
        }

        const LeaveMeetingBackend = async () => {
            notificationsWsRef.current?.close();
            await fetch(BackendAddressHttp + SetUrlParams(ApiEndpoints.leaveMeeting, meetingID), {
                method: 'POST',
                credentials: 'include',
            })
        }

        LeaveMeetingDasher();
        await LeaveMeetingBackend();
        LeaveMeetingFrontend(MeetingConfig.meetingState.left);
    }

    const KickFromMeeting = async (userID: string) => {
        await fetch(BackendAddressHttp + SetUrlParams(ApiEndpoints.kickParticipant, meetingID, userID), {
            method: 'POST',
            credentials: 'include',
        })
    }

    const GetExitText = () => {
        let title = '';
        switch (meetingState) {
            case MeetingConfig.meetingState.left:
                title = MeetingExitText.popUpTitles.left;
                break;
            case MeetingConfig.meetingState.kicked:
                title = MeetingExitText.popUpTitles.kicked;
                break;
            case MeetingConfig.meetingState.ended:
                title = MeetingExitText.popUpTitles.ended;
                break;
            case MeetingConfig.meetingState.wrongID:
                title = MeetingExitText.popUpTitles.wrongID;
                break;
            case MeetingConfig.meetingState.error:
                title = MeetingExitText.popUpTitles.error;
                break;
            default:
                title = MeetingExitText.popUpTitles.default;
                break;
        }
        return title
    }

    if (meetingState !== MeetingConfig.meetingState.active && meetingState !== MeetingConfig.meetingState.none) {

        const title = GetExitText();
        return (
            <OneButtonPopUp
                open={true}
                theme='success'
                title={title}
                buttonText={MeetingExitText.popUpButton}
                onButtonClick={() => {
                    goToHome();
                }}
            >
                {MeetingExitText.popUpSubtitle}
            </OneButtonPopUp>
        )
    }
    return (
        <div>
            <StyledMeetingGrid>
                {
                    <StyledMeetingGridTile>
                        <StyledVideo
                            ref={clientVideoRef}
                            autoPlay
                            playsInline
                            muted
                        />
                    </StyledMeetingGridTile>
                }
                {
                    participantsIDs.map((id) => (
                        <StyledMeetingGridTile key={id}>
                            <DashPlayer
                                userID={id}
                                url={SetUrlParams(DasherServerAddressHttp + ApiEndpoints.getStream, meetingID, id)}
                                menuOptions={hostOptions}
                            />
                        </StyledMeetingGridTile>
                    ))
                }
            </StyledMeetingGrid>

            <StyledMeetingFooter
                onLeaveMeeting={LeaveMeeting}
            />
        </div>
    )
}

export default MeetingPage