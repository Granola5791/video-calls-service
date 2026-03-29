import { TableRow, TableCell } from '@mui/material'
import { MeetingInfoText } from '../constants/hebrew-constants'
import type { MeetingInfo } from '../types/meetingInfo'
import { ClickableP } from '../styled-components/StyledText'

interface MeetingInfoRowProps {
    meetingInfo: MeetingInfo
    onTranscriptClick: (id: string) => void
    onSummaryClick: (id: string) => void
}

const MeetingInfoRow = ({ meetingInfo, onTranscriptClick, onSummaryClick }: MeetingInfoRowProps) => {
    return (
        <TableRow>
            <TableCell>{meetingInfo.id}</TableCell>
            <TableCell >{meetingInfo.name}</TableCell>
            <TableCell >{meetingInfo.date}</TableCell>
            <TableCell >{meetingInfo.hostName}</TableCell>
            <TableCell onClick={() => onTranscriptClick(meetingInfo.id)}><ClickableP>{MeetingInfoText.transcript}</ClickableP></TableCell>
            <TableCell onClick={() => onSummaryClick(meetingInfo.id)}><ClickableP>{MeetingInfoText.summary}</ClickableP></TableCell>
        </TableRow>
    )
}

export default MeetingInfoRow