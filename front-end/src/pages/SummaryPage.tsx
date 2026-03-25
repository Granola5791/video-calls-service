import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom';
import { SetUrlParams, UsersServer } from '../constants/backend-constants';
import { MeetingInfoText } from '../constants/hebrew-constants';
import { CenteredColumn } from '../styled-components/StyledBoxes';

const SplitLines = (text: string): string[] => {
    return text.split('\n');
}

const SummaryPage = () => {
    const { meetingID } = useParams();
    const [summaryLines, setSummaryLines] = useState<string[]>([]);

    useEffect(() => {
        const fetchSummary = async () => {
            const response = await fetch(SetUrlParams(UsersServer.httpAddress + UsersServer.api.getSummary, meetingID), {
                method: 'GET',
                credentials: 'include',
            });
            const data = await response.json() as string;
            const text = SplitLines(data);
            setSummaryLines(text);
        };
        fetchSummary();
    }, []);

    return (
        <CenteredColumn>
            <h1>{MeetingInfoText.summary}</h1>
            <div>
                {summaryLines.map((line, index) =>
                    <p key={index}>
                        {line.trim() || '\u00a0'}
                    </p>
                )}
            </div>
        </ CenteredColumn>
    )
}

export default SummaryPage