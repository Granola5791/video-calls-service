import { useEffect, useState } from 'react'
import { SetUrlParams, UsersServer } from '../constants/backend-constants';
import { CenteredColumn } from '../styled-components/StyledBoxes';
import { MeetingInfoText } from '../constants/hebrew-constants';
import type { UserTranscript, TranscriptSegment, Index2D } from '../types/transcripts';

interface MeetingTranscriptPageProps {
    meetingID: string;
}

const NormalizeTranscript = (transcript: string): TranscriptSegment[] => {
    const segments = transcript.split('\n');
    return Array.from(segments, NormalizeTranscriptSegment);
}

const NormalizeTranscriptSegment = (line: string): TranscriptSegment => {
    // 1. Find the end of the 'start' string
    const firstSpace = line.indexOf(' ');
    const start = parseFloat(line.slice(0, firstSpace));

    // 2. Find the end of the 'end' string
    const secondSpace = line.indexOf(' ', firstSpace + 1);
    const end = parseFloat(line.slice(firstSpace + 1, secondSpace));

    // 3. The rest is the text
    const text = line.slice(secondSpace + 1);

    return { text, start, end };
}

const GetIndices = (arr: UserTranscript[]) => arr.flatMap((row, y) => row.segments.map((_, x) => ({ x, y })));

const MeetingTranscriptPage = ({ meetingID }: MeetingTranscriptPageProps) => {
    const [transcript, setTranscript] = useState<UserTranscript[]>([]);
    const [orderedIndices, setOrderedIndices] = useState<Index2D[]>([]);

    useEffect(() => {
        const fetchTranscript = async (participantID: number): Promise<UserTranscript> => {
            const response = await fetch(SetUrlParams(UsersServer.httpAddress + UsersServer.api.getTranscript, meetingID, participantID), {
                method: 'GET',
                credentials: 'include',
            });
            const data = await response.json();
            const { transcript, username } = data;
            return { segments: NormalizeTranscript(transcript), username: username };
        };

        const fetchMeetingParticipantsIDs = async () => {
            const response = await fetch(SetUrlParams(UsersServer.httpAddress + UsersServer.api.getAllMeetingParticipants, meetingID), {
                method: 'GET',
                credentials: 'include',
            })
            const IDs = await response.json() as string[];
            let arr = [] as UserTranscript[];
            for (let i = 0; i < IDs.length; i++) {
                arr.push(await fetchTranscript(parseInt(IDs[i])));
            }
            setTranscript(arr);
            let indices = GetIndices(arr);
            indices.sort((a, b) => arr[a.y].segments[a.x].start - arr[b.y].segments[b.x].start);
            setOrderedIndices(indices);
        }

        fetchMeetingParticipantsIDs();
    }, []);

    return (
        <CenteredColumn>
            <h1>{MeetingInfoText.transcript}</h1>
            <div>
                {orderedIndices.map((index) => (
                    <span style={{ display: 'flex', flexDirection: 'row', gap: '5px' }}>
                        <b>{transcript[index.y].username}:</b>
                        <b style={{ color: 'green' }}>
                            [{transcript[index.y].segments[index.x].start}]
                        </b>
                        <div>
                            {transcript[index.y].segments[index.x].text}
                        </div>
                        <b style={{ color: 'red' }}>
                            [{transcript[index.y].segments[index.x].end}]
                        </b>
                    </span>
                ))}
            </div>
        </CenteredColumn>
    )
}

export default MeetingTranscriptPage