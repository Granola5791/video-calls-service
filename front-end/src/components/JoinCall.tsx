import React from 'react'
import OneButtonPopUp from './OneButtonPopUp'
import { StyledTextField } from '../styled-components/StyledTextFields'

interface JoinCallProps {
    onJoin?: () => void
}

const JoinCall = ({ onJoin}: JoinCallProps) => {
    return (
        <OneButtonPopUp
            open={true}
            title='Join Call'
            buttonText='Join'
            onButtonClick={onJoin}
        >
            <StyledTextField label='Meeting ID' />
        </OneButtonPopUp>
    )
}

export default JoinCall