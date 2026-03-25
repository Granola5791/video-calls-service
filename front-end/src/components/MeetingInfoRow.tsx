import { TableRow, TableCell } from '@mui/material'
import { MeetingInfoText } from '../constants/hebrew-constants'
import type { MeetingInfo } from '../types/meetingInfo'
import { ClickableP } from '../styled-components/StyledText'

interface MeetingInfoRowProps {
    meetingInfo: MeetingInfo
    onTranscriptClick: (id : string) => void
    onSummaryClick: (id : string) => void
}

const MeetingInfoRow = ({ meetingInfo, onTranscriptClick, onSummaryClick }: MeetingInfoRowProps) => {
    return (
        <TableRow>
            <TableCell align='right'>{meetingInfo.id}</TableCell>
            <TableCell align='right'>{meetingInfo.name}</TableCell>
            <TableCell align='right'>{meetingInfo.date}</TableCell>
            <TableCell align='right'>{meetingInfo.hostName}</TableCell>
            <TableCell align='right' onClick={() => onTranscriptClick(meetingInfo.id)}><ClickableP>{MeetingInfoText.transcript}</ClickableP></TableCell>
            <TableCell align='right' onClick={() => onSummaryClick(meetingInfo.id)}><ClickableP>{MeetingInfoText.summary}</ClickableP></TableCell>
        </TableRow>
    )
}

export default MeetingInfoRow