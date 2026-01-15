import React, { useEffect, useRef, useState } from 'react'
import WebCam from '../components/WebCam'
import { CenteredColumn } from '../styled-components/StyledBoxes'
import WebSocketWebCam from '../components/WebSocketWebCam'
import { ApiEndpoints, DasherServerAddressWS, DasherServerAddressHttp, BackendAddressHttp, BackendAddressWS, SetUrlParams, CallEventTypes } from '../constants/backend-constants'
import { LongButton, LongButtonFilled } from '../styled-components/StyledButtons'
import DashPlayer from '../components/DashPlayer'
import { StyledTextField } from '../styled-components/StyledTextFields'
import { Sleep } from '../utils/sleep'

const NormalizeMeetingIDs = (meetingIDs: unknown): string[] => {
    return Array.isArray(meetingIDs) ? meetingIDs : [];
}

const HomePage = () => {
    const wsRef = useRef<WebSocket | null>(null);
    const [isStreaming1, setIsStreaming1] = useState(false);
    const [isStreaming2, setIsStreaming2] = useState(false);
    const [meetingID, setMeetingID] = useState('');
    const [participantsIDs, setParticipantsIDs] = useState<string[]>([]);

    const CreateMeeting = async () => {
        const res1 = await fetch(BackendAddressHttp + ApiEndpoints.createMeeting, {
            method: 'POST',
            credentials: 'include',
        });
        if (!res1.ok) {
            console.error('Failed to create meeting in backend:', res1.statusText);
            return;
        }
        const meetingId = await res1.text();
        setMeetingID(meetingId);

        const res2 = await fetch(DasherServerAddressHttp + ApiEndpoints.createMeeting, {
            method: 'POST',
            credentials: 'include',
        });
        if (!res2.ok) {
            console.error('Failed to create meeting:', res2.statusText);
            return;
        }

        setIsStreaming1(true);
        JoinMeetingBackend(meetingId);

    }

    const JoinMeetingBackend = async (meetingID: string) => {
        let res = await fetch(BackendAddressHttp + SetUrlParams(ApiEndpoints.joinMeeting, meetingID), {
            method: 'POST',
            credentials: 'include',
        })
        if (res.ok) {
            const data = await res.json();
            setParticipantsIDs(NormalizeMeetingIDs(data));
        }

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
        }
    }

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
        await JoinMeetingBackend(meetingID);
        setIsStreaming1(true);
    }

    return (
        <div>
            <CenteredColumn>
                <LongButtonFilled onClick={CreateMeeting}>Create Meeting</LongButtonFilled>
                {isStreaming1 &&
                    <WebSocketWebCam wsUrl={DasherServerAddressWS + SetUrlParams(ApiEndpoints.startStream, meetingID)} />
                }
                {
                    participantsIDs.map((participantID) => (
                        <DashPlayer key={participantID} url={DasherServerAddressHttp + ApiEndpoints.getStream.prefix + meetingID + '/' + participantID + ApiEndpoints.getStream.suffix} />
                    ))
                }
                <StyledTextField value={meetingID} onChange={(e) => setMeetingID(e.target.value)} placeholder="Meeting ID" />
                <LongButton onClick={() => JoinMeeting(meetingID)}>Join Meeting</LongButton>
            </CenteredColumn>
        </div>
    )
}

export default HomePage 