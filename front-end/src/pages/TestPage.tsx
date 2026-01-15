import React from 'react'
import WebSocketWebCam from '../components/WebSocketWebCam'
import DashPlayer from '../components/DashPlayer'
import { DasherServerAddressWS, ApiEndpoints, DasherServerAddressHttp } from '../constants/backend-constants'

const TestPage = () => {
  return (
    <div>
        <DashPlayer url={DasherServerAddressHttp + ApiEndpoints.getStream + '/1'} />
        <DashPlayer url={DasherServerAddressHttp + ApiEndpoints.getStream + '/2'} />
    </div>
  )
}

export default TestPage