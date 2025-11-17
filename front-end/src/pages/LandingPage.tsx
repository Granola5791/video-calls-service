import React from 'react'
import { StyledTitle } from '../styled-components/StyledText'
import { CenteredColumn } from '../styled-components/StyledBoxes'
import { LongButton, LongButtonFilled } from '../styled-components/StyledButtons'
import { useNavigation } from '../utils/navigation'

const LandingPage = () => {

    const {
        goToSignUp: GoToSignUp,
        goToLogIn: GoToLogIn,
    } = useNavigation();

    return (
        <>
            <StyledTitle>Welcome to the Landing Page</StyledTitle>
            <CenteredColumn>
                <LongButtonFilled>Join a meeting</LongButtonFilled>
                <LongButton onClick={GoToSignUp}>sign up</LongButton>
                <LongButton onClick={GoToLogIn}>log in</LongButton>
            </CenteredColumn>
        </>
    )
}

export default LandingPage