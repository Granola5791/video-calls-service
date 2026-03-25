import { useNavigate } from "react-router-dom";
import { RouterPaths } from "../constants/general-contants";

export function useNavigation() {
    const navigate = useNavigate();

    return {
        goToSignUp: () => navigate(RouterPaths.signup),
        goToLogIn: () => navigate(RouterPaths.login),
        goToHome: () => navigate(RouterPaths.home),
        goToMeeting: (meetingID: string) => navigate(`/meeting/${meetingID}`),
        goToMeetingTranscript: (meetingID: string) => navigate(`/transcript/${meetingID}`),
        goToTranscripts: () => navigate(RouterPaths.meetingInfo),
        goToSummary: (meetingID: string) => navigate(`/summary/${meetingID}`),
    };
}