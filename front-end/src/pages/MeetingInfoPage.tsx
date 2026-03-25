import { AddQueryParams, UsersServer } from '../constants/backend-constants';
import { useNavigation } from '../utils/navigation';
import type { MeetingInfo } from '../types/meetingInfo';
import MeetingInfoTable from '../components/MeetingInfoTable';
import { CenteredColumn, CenteredRow } from '../styled-components/StyledBoxes';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { MeetingInfoText } from '../constants/hebrew-constants';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import { MediumButtonFilled } from '../styled-components/StyledButtons';
import dayjs, { Dayjs } from 'dayjs';
import type { QueryParam } from '../types/queryParam';
import { useState } from 'react';
import { StyledTextField } from '../styled-components/StyledTextFields';


const MeetingInfoPage = () => {
    const datePickerSlotProps = {
        textField: {
            variant: 'filled',
            InputLabelProps: {
                sx: {
                    width: '100%',
                    textAlign: 'right',
                },
            },
        },
    }
    const [meetings, setMeetings] = useState<MeetingInfo[]>([]);
    const [fromDate, setFromDate] = useState<Dayjs>(dayjs('2023-01-01T00:00:00.000Z'));
    const [toDate, setToDate] = useState<Dayjs>(dayjs(Date.now()));
    const [hostName, setHostName] = useState<string>('');
    const {
        goToMeetingTranscript,
        goToSummary,
    } = useNavigation();

    const fetchMeetings = async (fromDate: Dayjs, toDate: Dayjs): Promise<MeetingInfo[]> => {
        const queryParams = [
            {key: UsersServer.api.queryParams.from, value: fromDate.toISOString()},
            {key: UsersServer.api.queryParams.to, value: toDate.toISOString()},
            {key: UsersServer.api.queryParams.host_username, value: hostName},
        ] as QueryParam[]
        console.log(fromDate.toISOString(), toDate.toISOString());
        let url = UsersServer.httpAddress + UsersServer.api.getMeetingInfos
        url = AddQueryParams(url, queryParams)
        const response = await fetch(url, {
            method: 'GET',
            credentials: 'include',
        });
        const data = await response.json();
        const receivedMeetings = Array.from(data, (meeting: any) => ({ id: meeting.id, name: "placeholder", date: meeting.created_at, hostName: meeting.host_username } as MeetingInfo));
        return receivedMeetings;
    };

    const OnSearch = async () => {
        setMeetings(await fetchMeetings(fromDate, toDate));
    }

    return (
        <CenteredColumn>
            <h1>{MeetingInfoText.title}</h1>
            <StyledTextField
                label={MeetingInfoText.hostName}
                onChange={(e) => setHostName(e.target.value)}
            />
            <CenteredRow>
                <LocalizationProvider dateAdapter={AdapterDayjs}>
                    <DateTimePicker
                        label={MeetingInfoText.fromDate}
                        slotProps={datePickerSlotProps as any}
                        onChange={(value) => value?.isValid() && setFromDate(value)}
                        defaultValue={fromDate}
                    />
                    <DateTimePicker
                        label={MeetingInfoText.toDate}
                        slotProps={datePickerSlotProps as any}
                        onChange={(value) => value?.isValid() && setToDate(value)}
                        defaultValue={toDate}
                    />
                    <MediumButtonFilled onClick={OnSearch}>{MeetingInfoText.search}</MediumButtonFilled>
                </LocalizationProvider>
            </CenteredRow>
            {meetings.length > 0 && <MeetingInfoTable meetings={meetings} onTranscriptClick={goToMeetingTranscript} onSummaryClick={goToSummary} />}
        </CenteredColumn>
    )
}

export default MeetingInfoPage