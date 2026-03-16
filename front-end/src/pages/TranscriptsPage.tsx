import React, { useEffect } from 'react'
import { UsersServer } from '../constants/backend-constants';
import { useNavigation } from '../utils/navigation';

const TranscriptsPage = () => {
    const [meetings, setMeetings] = React.useState<string[]>([]);
    const {
        goToTranscript,
    } = useNavigation();

    useEffect(() => {
        const fetchMeetings = async () => {
            const response = await fetch(UsersServer.httpAddress + UsersServer.api.getTranscriptionMeetings);
            const data = await response.json();
            setMeetings(data);
        };
        fetchMeetings();
    }, []);

    return (
        <div>
            {meetings.map((meeting) => (<div onClick={() => goToTranscript(meeting)}>{meeting}</div>))}
        </div>
    )
}

export default TranscriptsPage