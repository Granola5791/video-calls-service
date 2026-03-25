import { TableCell, TableHead, TableRow } from '@mui/material'
import { MeetingInfoText } from '../constants/hebrew-constants'
import { useState } from 'react';
import { StyledTableSortLabel } from '../styled-components/StyledTable';
import type { MeetingInfo } from '../types/meetingInfo';

interface MeetingInfoHeadProps {
    colIDs: (keyof MeetingInfo)[],
    onSort: (id: keyof MeetingInfo, direction: 'asc' | 'desc') => void
}

const MeetingInfoHead = ({ colIDs, onSort }: MeetingInfoHeadProps) => {
    const [sortedColumn, setSortedColumn] = useState<string | null>(null);
    const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('asc');

    const handleSort = (id: keyof MeetingInfo) => {
        if (sortedColumn === id) {
            const newSortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
            setSortDirection(newSortDirection);
            onSort(id, newSortDirection);
        } else {
            setSortedColumn(id);
            setSortDirection('asc');
            onSort(id, 'asc');
        }
    };

    return (
        <TableHead>
            <TableRow>
                {colIDs.map((id) => <TableCell align='right' key={id}>
                    <StyledTableSortLabel
                        active={id === sortedColumn}
                        direction={sortDirection}
                        onClick={() => handleSort(id)}
                    >
                        {MeetingInfoText[id as keyof typeof MeetingInfoText]}
                    </StyledTableSortLabel>
                </TableCell>)}
            </TableRow>
        </TableHead>
    )
}

export default MeetingInfoHead