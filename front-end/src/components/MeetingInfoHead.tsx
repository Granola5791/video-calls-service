import { TableCell, TableHead, TableRow } from '@mui/material'
import { MeetingInfoText } from '../constants/hebrew-constants'

const MeetingInfoHead = () => {
    return (
        <TableHead>
            <TableRow>
                <TableCell align='right'>{MeetingInfoText.id}</TableCell>
                <TableCell align='right'>{MeetingInfoText.name}</TableCell>
                <TableCell align='right'>{MeetingInfoText.date}</TableCell>
                <TableCell align='right'>{MeetingInfoText.hostName}</TableCell>
                <TableCell align='right'>{MeetingInfoText.transcript}</TableCell>
                <TableCell align='right'>{MeetingInfoText.summary}</TableCell>
            </TableRow>
        </TableHead>
    )
}

export default MeetingInfoHead