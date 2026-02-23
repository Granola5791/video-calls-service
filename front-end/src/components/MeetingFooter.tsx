import { LeaveMeetingButton } from '../styled-components/StyledButtons';
import CallEndIcon from '@mui/icons-material/CallEnd';
import VisibilityIcon from '@mui/icons-material/Visibility';
import { WatchingIcon } from '../styled-components/StyledIcons';
import { MeetingID } from '../styled-components/StyledText';
import { MeetingFooterText } from '../constants/hebrew-constants';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import { IconButton } from '@mui/material';

interface MeetingFooterProps {
    onLeaveMeeting: () => void;
    dangerSignOn?: boolean;
    className?: string;
    meetingID?: string;
}

const MeetingFooter = ({ onLeaveMeeting, className, dangerSignOn, meetingID = "" }: MeetingFooterProps) => {
    return (
        <div className={className}>
            <LeaveMeetingButton onClick={onLeaveMeeting}>
                <CallEndIcon />
            </LeaveMeetingButton>
            <MeetingID>
                {MeetingFooterText.meetingID + meetingID}
                <IconButton color='inherit' onClick={() => {navigator.clipboard.writeText(meetingID);}}>
                    <ContentCopyIcon/>
                </IconButton>
            </MeetingID>
            {dangerSignOn && <WatchingIcon />}
        </div>
    )
}

export default MeetingFooter