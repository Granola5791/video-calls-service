import { useState } from 'react'
import { CenteredColumn, CenteredScreen } from '../styled-components/StyledBoxes'
import { ApiEndpoints, DasherServerAddressHttp, UsersServerAddressHttp } from '../constants/backend-constants'
import { LongButton, LongButtonFilled } from '../styled-components/StyledButtons'
import { StyledTextField } from '../styled-components/StyledTextFields'
import { useNavigation } from '../utils/navigation'
import { General, HomePageText, MenuOptions } from '../constants/hebrew-constants'
import TwoButtonPopUp from '../components/TwoButtonPopUp'
import { Colors } from '../constants/general-contants'
import { StyledTitle } from '../styled-components/StyledText'
import HamburgerMenu from '../components/HamburgerMenu'

const HomePage = () => {
    const [meetingID, setMeetingID] = useState('');
    const [openJoinMeetingPopUp, setOpenJoinMeetingPopUp] = useState(false);
    const {
        goToMeeting,
    } = useNavigation();

    const CreateMeeting = async () => {
        const res1 = await fetch(UsersServerAddressHttp + ApiEndpoints.createMeeting, {
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

    const LogOut = () => {
        const LogoutFrontend = () => {
            localStorage.clear();
            window.location.href = '/';
        }

        const LogoutBackend = async () => {
            await fetch(UsersServerAddressHttp + ApiEndpoints.logOut, {
                method: 'POST',
                credentials: 'include',
            });
        }

        LogoutBackend().then(LogoutFrontend);
    }

    return (
        <CenteredScreen>
            <HamburgerMenu
                topButtons={[{ text: MenuOptions.disconnect, onClick: LogOut }]}
            />
            <CenteredColumn>

                <StyledTitle>{General.appName}</StyledTitle>
                <CenteredColumn>
                    <LongButtonFilled onClick={CreateMeeting}>{HomePageText.createMeetingButton}</LongButtonFilled>
                    <LongButton onClick={() => setOpenJoinMeetingPopUp(true)}>{HomePageText.joinMeetingButton}</LongButton>

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