import React, { useEffect } from 'react'
import { useParams } from 'react-router-dom';
import { DasherServer, UsersServer, CallEventTypes, HttpStatusCodes, SetUrlParams } from '../constants/backend-constants';
import DashPlayer from '../components/DashPlayer';
import { StyledMeetingGrid } from '../styled-components/StyledBoxes';
import { StyledMeetingGridTile, StyledVideo } from '../styled-components/StyledVideos';
import OneButtonPopUp from '../components/OneButtonPopUp';
import { useNavigation } from '../utils/navigation';
import { MeetingConfig, StreamConfig } from '../constants/general-contants';
import { HostOptions, MeetingExitText } from '../constants/hebrew-constants';
import { StyledMeetingFooter } from '../styled-components/StyledFooters';
import { NormalizeArray } from '../utils/array';
import MeetingPreparationPage from './MeetingPreparationPage';

type Participant = {
    id: string;
    name: string;
}

const MeetingPage = () => {
    const { meetingID } = useParams();
    const streamWsRef = React.useRef<WebSocket | null>(null);
    const notificationsWsRef = React.useRef<WebSocket | null>(null);
    const clientVideoRef = React.useRef<HTMLVideoElement>(null);
    const toCleanUpRef = React.useRef<MediaStream[]>([]); // Store streams that need to be cleaned up
    const recorderRef = React.useRef<MediaRecorder | null>(null);
    const [participants, setParticipants] = React.useState<Participant[]>([]);
    const [meetingState, setMeetingState] = React.useState(MeetingConfig.meetingState.none);
    const keepAliveIntervalIDRef = React.useRef(0);
    const leaveMeetingTimeoutIDRef = React.useRef(0);
    const [dangerSignOn, setDangerSignOn] = React.useState(false);
    const [hostOptions, setHostOptions] = React.useState<{ label: string, onClick: (userID: string) => void }[]>([]);
    const {
        goToHome,
    } = useNavigation();

    useEffect(() => {
        if (!meetingID) {
            setMeetingState(MeetingConfig.meetingState.wrongID);
            return;
        }
        CheckIfAbleToJoinMeeting(meetingID);

        return () => {
            CloseNotificationsConnection();
            StopStream();
            if (keepAliveIntervalIDRef.current !== 0) {
                clearInterval(keepAliveIntervalIDRef.current);
            }
        };
    }, []);

    const CheckIfAbleToJoinMeeting = async (meetingID: string) => {
        const res = await fetch(UsersServer.httpAddress + SetUrlParams(UsersServer.api.isAbleToJoinMeeting, meetingID), {
            method: 'GET',
            credentials: 'include',
        });
        switch (res.status) {
            case HttpStatusCodes.OK:
                setMeetingState(MeetingConfig.meetingState.none);
                break;
            case HttpStatusCodes.NotFound:
                setMeetingState(MeetingConfig.meetingState.wrongID);
                break;
            case HttpStatusCodes.Unauthorized:
                setMeetingState(MeetingConfig.meetingState.banned);
                break;
            default:
                setMeetingState(MeetingConfig.meetingState.error);
                break;
        }
    }

    const EnterMeeting = async () => {
        StopStream();
        meetingID && JoinMeeting(meetingID);
    }

    const CloseNotificationsConnection = () => {
        if (notificationsWsRef.current) {
            notificationsWsRef.current.close();
        }
    }

    const SendKeepAlive = async () => {
        const res = await fetch(UsersServer.httpAddress + SetUrlParams(UsersServer.api.keepAlive, meetingID), {
            method: 'POST',
            credentials: 'include',
        })
        if (!res.ok) {
            if (res.status === HttpStatusCodes.Unauthorized) {
                LeaveMeetingFrontend(MeetingConfig.meetingState.kicked);
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
        toCleanUpRef.current.push(stream);

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
        let res = await fetch(UsersServer.httpAddress + SetUrlParams(UsersServer.api.joinMeeting, meetingID), {
            method: 'POST',
            credentials: 'include',
        })
        if (res.ok) {
            const data = await res.json();
            const participants: Participant[] = NormalizeArray(data.participants).map((participant: any) => ({
                id: participant.user_id,
                name: participant.username,
            }));
            const is_host = data.is_host
            setParticipants(participants);
            if (is_host) {
                ActivateHostOptions();
            }
        } else if (res.status === HttpStatusCodes.NotFound) {
            LeaveMeetingFrontend(MeetingConfig.meetingState.wrongID);
        } else if (res.status === HttpStatusCodes.Unauthorized) {
            LeaveMeetingFrontend(MeetingConfig.meetingState.banned);
        }
        if (!res.ok) {
            throw new Error(res.statusText);
        }
    }

    const SubscribeToMeetingUpdates = async (meetingID: string) => {
        const ws = new WebSocket(UsersServer.wsAddress + SetUrlParams(UsersServer.api.getCallNotifications, meetingID));
        notificationsWsRef.current = ws;
        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            const participantID = data.participant_id
            const participantName = data.participant_name
            const eventType = data.event
            const eventValue = data.value // usually zero

            switch (eventType) {
                case CallEventTypes.participantJoined:
                    const newParticipant = { id: participantID, name: participantName };
                    setParticipants((prev) =>
                        prev.some((p) => p.id === participantID)
                            ? prev
                            : [...prev, newParticipant]
                    );
                    break;
                case CallEventTypes.participantLeft:
                    setParticipants((prevParticipants) => prevParticipants.filter(({ id }) => id !== participantID));
                    break;
                case CallEventTypes.participantKickedByHost:
                    setDangerSignOn(true);
                    setTimeout(() => { setDangerSignOn(false) }, eventValue * 1000); // eventValue here is the time until participant is kicked for certain
                    break;
                case CallEventTypes.meetingEnded:
                    LeaveMeetingFrontend(MeetingConfig.meetingState.ended);
                    break;
            }
        };

        ws.onclose = () => {
            notificationsWsRef.current = null;
        };
    };

    const JoinMeetingDasher = async (meetingID: string) => {
        let res = await fetch(DasherServer.httpAddress + SetUrlParams(DasherServer.api.joinMeeting, meetingID), {
            method: 'POST',
            credentials: 'include',
        })
        if (!res.ok) {
            return;
        }
    }

    const JoinMeeting = async (meetingID: string) => {
        try {
            setMeetingState(MeetingConfig.meetingState.active);
            await JoinMeetingBackend(meetingID);
            await JoinMeetingDasher(meetingID);
            await StartStreaming(DasherServer.wsAddress + SetUrlParams(DasherServer.api.startStream, meetingID));
            await SubscribeToMeetingUpdates(meetingID);
            ContinuouslySendKeepAlive();
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
        toCleanUpRef.current.forEach(stream => stream.getTracks().forEach(track => track.stop()));
    };

    const LeaveMeetingFrontend = async (state: number) => {
        setMeetingState(state);
        StopStream();
        const timeOutId = setTimeout(() => { goToHome() }, MeetingConfig.exitWaitTimeMs)
        leaveMeetingTimeoutIDRef.current = timeOutId
    }

    const LeaveMeeting = async () => {
        const LeaveMeetingDasher = () => {
            notificationsWsRef.current?.close();
        }

        const LeaveMeetingBackend = async () => {
            notificationsWsRef.current?.close();
            await fetch(UsersServer.httpAddress + SetUrlParams(UsersServer.api.leaveMeeting, meetingID), {
                method: 'POST',
                credentials: 'include',
            })
        }

        LeaveMeetingDasher();
        await LeaveMeetingBackend();
        LeaveMeetingFrontend(MeetingConfig.meetingState.left);
    }

    const KickFromMeeting = async (userID: string) => {
        await fetch(UsersServer.httpAddress + SetUrlParams(UsersServer.api.kickParticipant, meetingID, userID), {
            method: 'POST',
            credentials: 'include',
        })
    }

    const ToggleMute = () => {
        const stream = clientVideoRef.current?.srcObject as MediaStream;

        if (stream) {
            stream.getAudioTracks().forEach((track) => {
                track.enabled = !track.enabled;
            });
        }
    };

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
            case MeetingConfig.meetingState.banned:
                title = MeetingExitText.popUpTitles.banned;
                break;
            default:
                title = MeetingExitText.popUpTitles.default;
                break;
        }
        return title
    }

    if (meetingState === MeetingConfig.meetingState.none) {
        return (
            <MeetingPreparationPage onEnterMeeting={EnterMeeting} />
        )
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
                    clearTimeout(leaveMeetingTimeoutIDRef.current);
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
                    participants.map((participant) => (
                        <StyledMeetingGridTile key={participant.id}>
                            <DashPlayer
                                userID={participant.id}
                                userName={participant.name}
                                url={SetUrlParams(DasherServer.httpAddress + DasherServer.api.getStream, meetingID, participant.id)}
                                menuOptions={hostOptions}
                            />
                        </StyledMeetingGridTile>
                    ))
                }
            </StyledMeetingGrid>

            <StyledMeetingFooter
                onLeaveMeeting={LeaveMeeting}
                dangerSignOn={dangerSignOn}
                meetingID={meetingID}
                toggleMuteFunc={ToggleMute}
            />
        </div>
    )
}

export default MeetingPage