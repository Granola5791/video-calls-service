import { Table, TableBody, TableContainer, TablePagination, TableRow } from '@mui/material'
import MeetingInfoHead from './MeetingInfoHead'
import type { MeetingInfo } from '../types/meetingInfo'
import MeetingInfoRow from './MeetingInfoRow'
import Paper from '@mui/material/Paper';
import { useEffect, useState } from 'react';
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
    const [displayRows, setDisplayRows] = useState<MeetingInfo[]>([]);
    const [sortedColumn, setSortedColumn] = useState<keyof MeetingInfo | null>(null);
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc');

    useEffect(() => {
        const newDisplay = meetings.sort((a, b) => {
            if (sortedColumn === null) { return 0; }

            let ret = 0;
            if (a[sortedColumn] < b[sortedColumn]) { ret = -1; }
            else if (a[sortedColumn] > b[sortedColumn]) { ret = 1; }

            if (sortDirection === 'desc') { ret *= -1; }
            return ret;
        }).slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);
        setDisplayRows(newDisplay);
    }, [page, rowsPerPage, sortDirection, sortedColumn, meetings]);

    const onSort = (id: keyof MeetingInfo, direction: 'asc' | 'desc') => {
        setSortedColumn(id);
        setSortDirection(direction);
    }

    return (
        <TableContainer component={Paper}>
            <Table>
                <MeetingInfoHead colIDs={MeetingInfoGeneral.colIDs as (keyof MeetingInfo)[]} onSort={onSort} />
                <TableBody>
                    {displayRows.map((meeting) => (
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
                            onRowsPerPageChange={(event) => { setRowsPerPage(parseInt(event.target.value, 10)); setPage(0) }}
                        />
                    </TableRow>
                </StyledTableFooter>
            </Table>
        </TableContainer>
    )
}

export default MeetingInfoTable