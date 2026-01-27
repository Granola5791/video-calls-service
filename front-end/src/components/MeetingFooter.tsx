import { LeaveMeetingButton } from '../styled-components/StyledButtons';
import CallEndIcon from '@mui/icons-material/CallEnd';

interface MeetingFooterProps {
    onLeaveMeeting: () => void;
    className?: string;
}

const MeetingFooter = ({ onLeaveMeeting, className }: MeetingFooterProps) => {
    return (
        <div className={className}>
            <LeaveMeetingButton onClick={onLeaveMeeting}>
                <CallEndIcon />
            </LeaveMeetingButton>
        </div>
    )
}

export default MeetingFooter