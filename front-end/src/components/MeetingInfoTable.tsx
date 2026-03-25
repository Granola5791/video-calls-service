import { Table, TableBody, TableContainer, TablePagination, TableRow } from '@mui/material'
import MeetingInfoHead from './MeetingInfoHead'
import type { MeetingInfo } from '../types/meetingInfo'
import MeetingInfoRow from './MeetingInfoRow'
import Paper from '@mui/material/Paper';
import { useState } from 'react';
import { MeetingInfoGeneral } from '../constants/general-contants';
import { StyledTableFooter } from '../styled-components/StyledTable';

interface MeetingInfoTableProps {
    meetings: MeetingInfo[]
    onTranscriptClick: (id: string) => void
    onSummaryClick: (id: string) => void
}

const MeetingInfoTable = ({ meetings, onTranscriptClick, onSummaryClick }: MeetingInfoTableProps) => {
    const [rowsPerPage, setRowsPerPage] = useState(5);
    const [page, setPage] = useState(0);

    return (
        <TableContainer component={Paper}>
            <Table>
                <MeetingInfoHead />
                <TableBody>
                    {meetings.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage).map((meeting) => (
                        <MeetingInfoRow
                            key={meeting.id}
                            meetingInfo={meeting}
                            onTranscriptClick={onTranscriptClick}
                            onSummaryClick={onSummaryClick}
                        />
                    ))}
                </TableBody>
                <StyledTableFooter>
                    <TableRow>
                        <TablePagination
                            rowsPerPageOptions={MeetingInfoGeneral.rowsPerPageOptions}
                            rowsPerPage={rowsPerPage}
                            count={meetings.length}
                            page={page}
                            onPageChange={(event, newPage) => setPage(newPage)}
                            onRowsPerPageChange={(event) => {setRowsPerPage(parseInt(event.target.value, 10)); setPage(0)}}
                        />
                    </TableRow>
                </StyledTableFooter>
            </Table>
        </TableContainer>
    )
}

export default MeetingInfoTable