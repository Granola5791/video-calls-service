import { useNavigate } from "react-router-dom";

export function useNavigation() {
    const navigate = useNavigate();

    return {
        goToSignUp: () => navigate("/signup"),
        goToLogIn: () => navigate("/login"),
        goToJoinMeeting: () => navigate("/join-meeting"),
    };
}