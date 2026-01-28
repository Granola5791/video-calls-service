import React from 'react'
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import { CenteredColumn } from '../styled-components/StyledBoxes';
import { AdaptiveButtonFilled } from '../styled-components/StyledButtons';

interface OneButtonPopUpProps {
    open: boolean
    theme?: string
    title?: string
    buttonText?: string
    onButtonClick?: () => void
    disabled?: boolean
    children?: React.ReactNode
}

const OneButtonPopUp = ({ open, theme = '', title = '', buttonText = 'OK', onButtonClick, disabled = false, children }: OneButtonPopUpProps) => {

    return (
        <Dialog open={open} onClose={() => { }}>
            <CenteredColumn>
                {title && <DialogTitle>{title}</DialogTitle>}
                {children}
                <AdaptiveButtonFilled
                    sx={{ color: theme }}
                    onClick={onButtonClick}
                    disabled={disabled}
                >
                    {buttonText}
                </AdaptiveButtonFilled>
            </CenteredColumn>
        </Dialog>
    )
}

export default OneButtonPopUp