import OneButtonPopUp from "../components/OneButtonPopUp";
import { MeetingConfig } from "../constants/general-contants";
import { MeetingExitText } from "../constants/hebrew-constants";

const GetExitText = (meetingState: number) => {
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

interface MeetingEndPageProps {
    meetingState: number;
    onExit: () => void;
}

const MeetingEndPage = ({ meetingState, onExit }: MeetingEndPageProps) => {
    const title = GetExitText(meetingState);

    return (
        <OneButtonPopUp
            open={true}
            theme='success'
            title={title}
            buttonText={MeetingExitText.popUpButton}
            onButtonClick={onExit}
        >
            {MeetingExitText.popUpSubtitle}
        </OneButtonPopUp>
    )
}

export default MeetingEndPage