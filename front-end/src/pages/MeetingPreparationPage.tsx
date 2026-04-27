import { useEffect, useRef, useState } from "react"
import { StartMeetingText } from "../constants/hebrew-constants"
import { CenteredColumn } from "../styled-components/StyledBoxes"
import { LongButtonFilled } from "../styled-components/StyledButtons"
import { StyledTitle } from "../styled-components/StyledText"

interface MeetingPreparationPageProps {
    onEnterMeeting: () => void
}

const MeetingPreparationPage = ({ onEnterMeeting }: MeetingPreparationPageProps) => {
    const previewVideoRef = useRef<HTMLVideoElement>(null);
    const [previewOn, setPreviewOn] = useState(false);
    const toCleanUpRef = useRef<MediaStream[]>([]);

    useEffect(() => {
        ShowPreview();
        return () => {
            toCleanUpRef.current.forEach(stream => stream.getTracks().forEach(track => track.stop()));
        };
    }, []);

    const ShowPreview = async () => {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
            toCleanUpRef.current.push(stream);
            if (previewVideoRef.current) {
                previewVideoRef.current.srcObject = stream;
                await previewVideoRef.current.play();
            }
            setPreviewOn(true);
        } catch (error) {
            setPreviewOn(false);
        }
    }

    return (
        <CenteredColumn>
            <StyledTitle>{StartMeetingText.title}</StyledTitle>
            {!previewOn && <h4>{StartMeetingText.allowCameraAccess}</h4>}
            <video
                ref={previewVideoRef}
                autoPlay
                playsInline
                muted
                width="40%"
            />
            <LongButtonFilled
                onClick={onEnterMeeting}
                disabled={!previewOn}
            >
                {StartMeetingText.enterMeetingButton}
            </LongButtonFilled>
        </CenteredColumn>
    )
}

export default MeetingPreparationPage