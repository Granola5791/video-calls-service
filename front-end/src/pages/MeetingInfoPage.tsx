import { AddQueryParams, UsersServer } from '../constants/backend-constants';
import type { MeetingInfo } from '../types/meetingInfo';
import MeetingInfoTable from '../components/MeetingInfoTable';
import { CenteredColumn, CenteredRow } from '../styled-components/StyledBoxes';
import { MeetingInfoText } from '../constants/hebrew-constants';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import { MediumButtonFilled } from '../styled-components/StyledButtons';
import dayjs, { Dayjs } from 'dayjs';
import type { QueryParam } from '../types/queryParam';
import { useState } from 'react';
import { StyledTextField } from '../styled-components/StyledTextFields';
import OneButtonPopUp from '../components/OneButtonPopUp';
import MeetingTranscriptPage from './MeetingTranscriptPage';
import SummaryPage from './SummaryPage';

const MeetingInfoPage = () => {
    const [meetings, setMeetings] = useState<MeetingInfo[]>([]);
    const [fromDate, setFromDate] = useState<Dayjs>(dayjs('2023-01-01T00:00:00.000Z'));
    const [toDate, setToDate] = useState<Dayjs>(dayjs(Date.now()));
    const [hostName, setHostName] = useState<string>('');
    const [meetingName, setMeetingName] = useState<string>('');
    const [pickedMeetingId, setPickedMeetingId] = useState<string | null>(null);
    const [whatToDisplay, setWhatToDisplay] = useState<'summary' | 'transcript' | null>(null);

    const fetchMeetings = async (fromDate: Dayjs, toDate: Dayjs, hostName: string, meetingName: string): Promise<MeetingInfo[]> => {
        const queryParams = [
            { key: UsersServer.api.queryParams.from, value: fromDate.toISOString() },
            { key: UsersServer.api.queryParams.to, value: toDate.toISOString() },
            { key: UsersServer.api.queryParams.hostUsername, value: hostName },
            { key: UsersServer.api.queryParams.meetingName, value: meetingName },
        ] as QueryParam[]
        let url = UsersServer.httpAddress + UsersServer.api.getMeetingInfos
        url = AddQueryParams(url, queryParams)
        const response = await fetch(url, {
            method: 'GET',
            credentials: 'include',
        });
        const data = await response.json();
        const receivedMeetings = Array.from(data, (meeting: any) => ({
            id: meeting.id,
            name: meeting.name,
            date: meeting.created_at,
            hostName: meeting.host_username,
        } as MeetingInfo));
        return receivedMeetings;
    };

    const OnSearch = async () => {
        setMeetings(await fetchMeetings(fromDate, toDate, hostName, meetingName));
    }

    return (
        <CenteredColumn>
            <h1>{MeetingInfoText.title}</h1>
            <CenteredRow>
                <StyledTextField
                    label={MeetingInfoText.hostName}
                    onChange={(e) => setHostName(e.target.value)}
                />
                <StyledTextField
                    label={MeetingInfoText.meetingName}
                    onChange={(e) => setMeetingName(e.target.value)}
                />
            </CenteredRow>
            <CenteredRow>
                <DateTimePicker
                    label={MeetingInfoText.fromDate}
                    onChange={(value) => value?.isValid() && setFromDate(value)}
                    defaultValue={fromDate}
                    localeText={MeetingInfoText.hebrewPickerLabels}
                />
                <DateTimePicker
                    label={MeetingInfoText.toDate}
                    onChange={(value) => value?.isValid() && setToDate(value)}
                    defaultValue={toDate}
                    localeText={MeetingInfoText.hebrewPickerLabels}
                />
                <MediumButtonFilled onClick={OnSearch}>{MeetingInfoText.search}</MediumButtonFilled>
            </CenteredRow>
            {meetings.length > 0 && <MeetingInfoTable
                meetings={meetings}
                onTranscriptClick={(meetingID) => { setPickedMeetingId(meetingID); setWhatToDisplay('transcript'); }}
                onSummaryClick={(meetingID) => { setPickedMeetingId(meetingID); setWhatToDisplay('summary'); }}
            />}
            <OneButtonPopUp
                open={pickedMeetingId !== null}
                onButtonClick={() => setPickedMeetingId(null)}
            >
                {(() => {
                    switch (whatToDisplay) {
                        case 'summary':
                            return <SummaryPage meetingID={pickedMeetingId ?? ''} />;
                        case 'transcript':
                            return <MeetingTranscriptPage meetingID={pickedMeetingId ?? ''} />;
                        default:
                            return null;
                    }
                })()}
            </OneButtonPopUp>
        </CenteredColumn>
    )
}

export default MeetingInfoPage