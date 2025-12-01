import React from 'react'
import JoinCall from '../components/JoinCall'
import { BackendAddress, ApiEndpoints, HttpStatuses, HttpStatusCodes } from '../constants/backend-constants';
import { useNavigation } from '../utils/navigation';

const CheckMeetingLoader = async (meetingID: string) => {
    const res = await fetch(BackendAddress + ApiEndpoints.checkMeeting + meetingID, {
        method: 'GET',
        credentials: 'include',
    });
    return res.ok
}


const Test = () => {

    const {
        goToHome: GoToHome,
        goToMeeting: GoToMeeting,
    } = useNavigation();

    const OnJoin = async (meetingID: string) => {
        if (await CheckMeetingLoader(meetingID)) {
            GoToMeeting(meetingID);
        }
    }

    return (
        <div>
            <JoinCall />
        </div>
    )
}

export default Test