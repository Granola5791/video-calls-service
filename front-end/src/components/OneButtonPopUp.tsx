import React from 'react'
import Dialog from '@mui/material/Dialog';
import Button from '@mui/material/Button';
import DialogTitle from '@mui/material/DialogTitle';
import { CenteredColumn } from '../styled-components/StyledBoxes';

interface OneButtonPopUpProps {
    open: boolean
    theme?: 'error' | 'info' | 'success' | 'warning'
    title?: string
    buttonText?: string
    onButtonClick?: () => void
    children?: React.ReactNode
}

const OneButtonPopUp = ({ open, theme = 'success', title = '', buttonText = 'OK', onButtonClick, children }: OneButtonPopUpProps) => {

    return (
        <Dialog color='error' open={open} onClose={() => { }}>
            <CenteredColumn>
                {title && <DialogTitle>{title}</DialogTitle>}
                {children}
                <Button
                    size='small'
                    variant="contained"
                    color={theme}
                    onClick={onButtonClick}
                >
                    {buttonText}
                </Button>
            </CenteredColumn>
        </Dialog>
    )
}

export default OneButtonPopUp