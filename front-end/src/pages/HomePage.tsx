import React from 'react'
import WebCam from '../components/WebCam'
import { CenteredColumn } from '../styled-components/StyledBoxes'

const HomePage = () => {
    return (
        <div>
            <CenteredColumn>
                <WebCam />
            </CenteredColumn>
        </div>
    )
}

export default HomePage 