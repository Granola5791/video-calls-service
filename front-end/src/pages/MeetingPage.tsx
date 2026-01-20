import React, { useEffect } from 'react'
import { useParams } from 'react-router-dom';
import { ApiEndpoints, BackendAddressHttp, BackendAddressWS, CallEventTypes, DasherServerAddressHttp, DasherServerAddressWS, SetUrlParams } from '../constants/backend-constants';
import WebSocketWebCam from '../components/WebSocketWebCam';
import DashPlayer from '../components/DashPlayer';
import { StyledMeetingGrid } from '../styled-components/StyledBoxes';
import { StyledMeetingGridTile } from '../styled-components/StyledVideos';


const NormalizeMeetingIDs = (meetingIDs: unknown): string[] => {
    return Array.isArray(meetingIDs) ? meetingIDs : [];
}

const MeetingPage = () => {

    const { meetingID } = useParams();
    const wsRef = React.useRef<WebSocket | null>(null);
    const [participantsIDs, setParticipantsIDs] = React.useState<string[]>([]);
    const [isStreaming, setIsStreaming] = React.useState(false);

    useEffect(() => {
        if (meetingID) {
            JoinMeeting(meetingID);
        }
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
        await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
        await JoinMeetingDasher(meetingID);
        StartStreaming();
        await JoinMeetingBackend(meetingID);
        await SubscribeToMeetingUpdates(meetingID);
    }

    return (
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
    )
}

export default MeetingPage