import React from 'react'
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import { CenteredColumn, CenteredRow } from '../styled-components/StyledBoxes';
import { AdaptiveButton, AdaptiveButtonFilled } from '../styled-components/StyledButtons';

interface OneButtonPopUpProps {
    open: boolean
    buttonColor1?: string
    buttonColor2?: string
    title?: string
    buttonText1?: string
    buttonText2?: string
    button1Disabled?: boolean
    button2Disabled?: boolean
    onButtonClick1?: () => void
    onButtonClick2?: () => void
    children?: React.ReactNode
}

const TwoButtonPopUp = ({ open, buttonColor1 = '', buttonColor2 = '', title = '', buttonText1 = 'OK', buttonText2 = 'Cancel', button1Disabled = false, button2Disabled = false, onButtonClick1, onButtonClick2, children }: OneButtonPopUpProps) => {

    return (
        <Dialog open={open} onClose={() => { }}>
            <CenteredColumn>
                {title && <DialogTitle>{title}</DialogTitle>}
                {children}
                <CenteredRow>
                    <AdaptiveButtonFilled
                        sx={{ backgroundColor: buttonColor1 }}
                        onClick={onButtonClick1}
                        disabled={button1Disabled}
                    >
                        {buttonText1}
                    </AdaptiveButtonFilled>
                    <AdaptiveButton
                        sx={{ backgroundColor: buttonColor2 }}
                        onClick={onButtonClick2}
                        disabled={button2Disabled}
                    >
                        {buttonText2}
                    </AdaptiveButton>
                </CenteredRow>
            </CenteredColumn>
        </Dialog>
    )
}

export default TwoButtonPopUp