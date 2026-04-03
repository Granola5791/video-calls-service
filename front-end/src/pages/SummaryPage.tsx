import { useEffect, useState } from 'react'
import { SetUrlParams, UsersServer } from '../constants/backend-constants';
import { MeetingInfoText } from '../constants/hebrew-constants';
import { CenteredColumn, CenteredRow } from '../styled-components/StyledBoxes';
import DownloadButton from '../components/DownloadButton';

interface SummaryPageProps {
    meetingID: string
}

const SplitLines = (text: string): string[] => {
    return text.split('\n');
}

const SummaryPage = ({ meetingID }: SummaryPageProps) => {
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
            <CenteredRow>
                <h1>{MeetingInfoText.summary}</h1>
                <DownloadButton fileName="summary.txt" putTextOnClick={() => summaryLines.join('\n')} />
            </CenteredRow>
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