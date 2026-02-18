import { LeaveMeetingButton } from '../styled-components/StyledButtons';
import CallEndIcon from '@mui/icons-material/CallEnd';
import VisibilityIcon from '@mui/icons-material/Visibility';
import { WatchingIcon } from '../styled-components/StyledIcons';

interface MeetingFooterProps {
    onLeaveMeeting: () => void;
    dangerSignOn?: boolean;
    className?: string;
}

const MeetingFooter = ({ onLeaveMeeting, className, dangerSignOn }: MeetingFooterProps) => {
    return (
        <div className={className}>
            <LeaveMeetingButton onClick={onLeaveMeeting}>
                <CallEndIcon />
            </LeaveMeetingButton>
            {dangerSignOn && <WatchingIcon />}
        </div>
    )
}

export default MeetingFooter