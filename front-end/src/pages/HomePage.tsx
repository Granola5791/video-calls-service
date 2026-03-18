import { useState } from 'react'
import { CenteredColumn, CenteredScreen } from '../styled-components/StyledBoxes'
import { UsersServer, DasherServer, SetUrlParams } from '../constants/backend-constants'
import { LongButton, LongButtonFilled } from '../styled-components/StyledButtons'
import { StyledTextField } from '../styled-components/StyledTextFields'
import { useNavigation } from '../utils/navigation'
import { General, HomePageText } from '../constants/hebrew-constants'
import TwoButtonPopUp from '../components/TwoButtonPopUp'
import { Colors } from '../constants/general-contants'
import { StyledTitle } from '../styled-components/StyledText'
import Checkbox from '@mui/material/Checkbox';
import { FormControlLabel } from '@mui/material'

const HomePage = () => {
    const [meetingID, setMeetingID] = useState('');
    const [openJoinMeetingPopUp, setOpenJoinMeetingPopUp] = useState(false);
    const [openCreateMeetingPopUp, setOpenCreateMeetingPopUp] = useState(false);
    const [requireFace, setRequireFace] = useState(false);
    const {
        goToMeeting,
    } = useNavigation();

    const CreateMeeting = async () => {
        const res1 = await fetch(SetUrlParams(UsersServer.httpAddress + UsersServer.api.createMeeting, requireFace), {
            method: 'POST',
            credentials: 'include',
        });
        if (!res1.ok) {
            console.error('Failed to create meeting in backend:', res1.statusText);
            return;
        }
        const meetingId = await res1.text();
        setMeetingID(meetingId);

        const res2 = await fetch(DasherServer.httpAddress + DasherServer.api.createMeeting, {
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
        <CenteredScreen>
            <CenteredColumn>

                <StyledTitle>{General.appName}</StyledTitle>
                <CenteredColumn>
                    <LongButtonFilled onClick={() => setOpenCreateMeetingPopUp(true)}>{HomePageText.createMeetingButton}</LongButtonFilled>
                    <LongButton onClick={() => setOpenJoinMeetingPopUp(true)}>{HomePageText.joinMeetingButton}</LongButton>

                    <TwoButtonPopUp
                        open={openCreateMeetingPopUp}
                        onButtonClick1={CreateMeeting}
                        onButtonClick2={() => { setOpenCreateMeetingPopUp(false) }}
                        buttonText1={HomePageText.submitMeetingIDButton}
                        buttonText2={HomePageText.cancelMeetingIDButton}
                        buttonColor1={Colors.primary}
                    >
                            <FormControlLabel
                                label={HomePageText.requireFaceLabel}
                                control={<Checkbox onChange={() => {setRequireFace(!requireFace)}}/>}
                            />
                    </TwoButtonPopUp>

                    <TwoButtonPopUp
                        open={openJoinMeetingPopUp}
                        onButtonClick1={() => { goToMeeting(meetingID) }}
                        onButtonClick2={() => { setOpenJoinMeetingPopUp(false) }}
                        buttonText1={HomePageText.submitMeetingIDButton}
                        buttonText2={HomePageText.cancelMeetingIDButton}
                        buttonColor1={Colors.primary}
                        button1Disabled={!meetingID}
                    >
                        <StyledTextField value={meetingID} onChange={(e) => setMeetingID(e.target.value)} placeholder={HomePageText.meetingIDInputPlaceholder} />
                    </TwoButtonPopUp>

                </CenteredColumn>
            </CenteredColumn>
        </CenteredScreen>
    )
}

export default HomePage