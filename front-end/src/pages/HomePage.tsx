import React from 'react'
import WebCam from '../components/WebCam'
import { CenteredColumn } from '../styled-components/StyledBoxes'
import WebSocketWebCam from '../components/WebSocketWebCam'
import { ApiEndpoints, DasherServerAddressWS } from '../constants/backend-constants'

const HomePage = () => {
    return (
        <div>
            <CenteredColumn>
                <WebSocketWebCam wsUrl={DasherServerAddressWS + ApiEndpoints.startStream} />
            </CenteredColumn>
        </div>
    )
}

export default HomePage 