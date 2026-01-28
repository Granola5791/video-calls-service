import { StyledMeetingFooter } from '../styled-components/StyledFooters'

const TestPage = () => {
  return (
    <div>
        <StyledMeetingFooter onLeaveMeeting={() => {console.log("Leaving meeting")}} />
    </div>
  )
}

export default TestPage