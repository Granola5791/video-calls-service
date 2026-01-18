import { useState } from 'react'
import { CenteredColumn } from '../styled-components/StyledBoxes'
import { ApiEndpoints, DasherServerAddressHttp, BackendAddressHttp } from '../constants/backend-constants'
import { LongButton, LongButtonFilled } from '../styled-components/StyledButtons'
import { StyledTextField } from '../styled-components/StyledTextFields'
import OneButtonPopUp from '../components/OneButtonPopUp'
import { useNavigation } from '../utils/navigation'

const HomePage = () => {
    const [meetingID, setMeetingID] = useState('');
    const [openJoinMeetingPopUp, setOpenJoinMeetingPopUp] = useState(false);
    const {
        goToMeeting,
    } = useNavigation();

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

        goToMeeting(meetingId);
    }



    return (
        <div>
            <CenteredColumn>
                <LongButtonFilled onClick={() => setOpenJoinMeetingPopUp(true)}>Join Meeting</LongButtonFilled>
                <LongButton onClick={CreateMeeting}>Create Meeting</LongButton>

                <OneButtonPopUp
                    open={openJoinMeetingPopUp}
                    onButtonClick={() => goToMeeting(meetingID)}
                >
                    <StyledTextField value={meetingID} onChange={(e) => setMeetingID(e.target.value)} placeholder="Meeting ID" />
                </OneButtonPopUp>

            </CenteredColumn>
        </div>
    )
}

export default HomePage