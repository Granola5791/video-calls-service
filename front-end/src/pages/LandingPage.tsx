import React from 'react'
import { StyledTitle } from '../styled-components/StyledText'
import { CenteredColumn, CenteredFilledScreen } from '../styled-components/StyledBoxes'
import { LongButton, LongButtonFilled } from '../styled-components/StyledButtons'
import { useNavigation } from '../utils/navigation'
import { General, LandingPageText } from '../constants/hebrew-constants'

const LandingPage = () => {

    const {
        goToSignUp: GoToSignUp,
        goToLogIn: GoToLogIn,
    } = useNavigation();

    return (
        <CenteredFilledScreen>
            <CenteredColumn >
                <StyledTitle>{General.appName}</StyledTitle>
                <CenteredColumn>
                    <LongButtonFilled>{LandingPageText.joinMeetingButton}</LongButtonFilled>
                    <LongButton onClick={GoToSignUp}>{LandingPageText.signUpButton}</LongButton>
                    <LongButton onClick={GoToLogIn}>{LandingPageText.logInButton}</LongButton>
                </CenteredColumn>
            </CenteredColumn>
        </CenteredFilledScreen>
    )
}

export default LandingPage