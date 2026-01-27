import React, { useEffect } from 'react'
import { useParams } from 'react-router-dom';
import { ApiEndpoints, BackendAddressHttp, BackendAddressWS, CallEventTypes, DasherServerAddressHttp, DasherServerAddressWS, SetUrlParams } from '../constants/backend-constants';
import WebSocketWebCam from '../components/WebSocketWebCam';
import DashPlayer from '../components/DashPlayer';
import { StyledMeetingGrid } from '../styled-components/StyledBoxes';
import { StyledMeetingGridTile } from '../styled-components/StyledVideos';
import OneButtonPopUp from '../components/OneButtonPopUp';
import { useNavigation } from '../utils/navigation';
import { MeetingConfig } from '../constants/general-contants';
import { MeetingExitText } from '../constants/hebrew-constants';
import { StyledMeetingFooter } from '../styled-components/StyledFooters';


const NormalizeMeetingIDs = (meetingIDs: unknown): string[] => {
    return Array.isArray(meetingIDs) ? meetingIDs : [];
}

const MeetingPage = () => {

    const { meetingID } = useParams();
    const wsRef = React.useRef<WebSocket | null>(null);
    const streamRef = React.useRef<MediaStream | null>(null);
    const [participantsIDs, setParticipantsIDs] = React.useState<string[]>([]);
    const [isStreaming, setIsStreaming] = React.useState(false);
    const [meetingState, setMeetingState] = React.useState(MeetingConfig.meetingState.none);
    const {
        goToHome,
    } = useNavigation();

    useEffect(() => {
        if (meetingID) {
            JoinMeeting(meetingID);
        }

        return () => {
            if (wsRef.current) {
                wsRef.current.close();
            }
            CloseWebCam();
        };
    }, []);

    const StartStreaming = () => {
        setIsStreaming(true);
    };

    const JoinMeetingBackend = async (meetingID: string) => {
        let res = await fetch(BackendAddressHttp + SetUrlParams(ApiEndpoints.joinMeeting, meetingID), {
            method: 'POST',
            credentials: 'include',
        })
        if (res.ok) {
            const data = await res.json();
            setParticipantsIDs(NormalizeMeetingIDs(data));
        }
    }

    const SubscribeToMeetingUpdates = async (meetingID: string) => {
        const ws = new WebSocket(BackendAddressWS + SetUrlParams(ApiEndpoints.getCallNotifications, meetingID));
        wsRef.current = ws;
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
            wsRef.current = null;
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
        const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
        streamRef.current = stream;
        await JoinMeetingDasher(meetingID);
        StartStreaming();
        await JoinMeetingBackend(meetingID);
        await SubscribeToMeetingUpdates(meetingID);
        setMeetingState(MeetingConfig.meetingState.active);
    }

    const CloseWebCam = () => {
        streamRef.current?.getTracks().forEach((track) => track.stop());
        streamRef.current = null;
    }

    const LeaveMeeting = async () => {
        const LeaveMeetingFrontend = async (state: number) => {
            setMeetingState(state);
            CloseWebCam();
            setTimeout(() => { goToHome() }, MeetingConfig.exitWaitTimeMs)
        }

        const LeaveMeetingDasher = () => {
            wsRef.current?.close();
        }

        const LeaveMeetingBackend = async () => {
            await fetch(BackendAddressHttp + SetUrlParams(ApiEndpoints.leaveMeeting, meetingID), {
                method: 'POST',
                credentials: 'include',
            })
        }

        LeaveMeetingDasher();
        await LeaveMeetingBackend();
        LeaveMeetingFrontend(MeetingConfig.meetingState.left);
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
                        {isStreaming && <WebSocketWebCam
                            wsUrl={DasherServerAddressWS + SetUrlParams(ApiEndpoints.startStream, meetingID!)}
                        />
                        }
                    </StyledMeetingGridTile>
                }
                {
                    participantsIDs.map((id) => (
                        <StyledMeetingGridTile key={id}>
                            <DashPlayer url={SetUrlParams(DasherServerAddressHttp + ApiEndpoints.getStream, meetingID, id)} />
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