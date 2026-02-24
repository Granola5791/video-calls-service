import { LeaveMeetingButton } from '../styled-components/StyledButtons';
import CallEndIcon from '@mui/icons-material/CallEnd';
import { WatchingIcon, MutedIcon, UnmutedIcon } from '../styled-components/StyledIcons';
import { MeetingID } from '../styled-components/StyledText';
import { MeetingFooterText } from '../constants/hebrew-constants';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import { IconButton } from '@mui/material';
import { useState } from 'react';

interface MeetingFooterProps {
    onLeaveMeeting: () => void;
    dangerSignOn?: boolean;
    className?: string;
    meetingID?: string;
    toggleMuteFunc?: () => void;
}

const MeetingFooter = ({ onLeaveMeeting, className, dangerSignOn, meetingID = "", toggleMuteFunc: toggleMuteFunc = () => {} }: MeetingFooterProps) => {
    const [muted, setMuted] = useState(false);

    const OnToggleMute = () => {
        setMuted(!muted);
        toggleMuteFunc();
    }

    return (
        <div className={className}>
            <LeaveMeetingButton onClick={onLeaveMeeting}>
                <CallEndIcon />
            </LeaveMeetingButton>
            <MeetingID>
                {MeetingFooterText.meetingID + meetingID}
                <IconButton color='inherit' onClick={() => { navigator.clipboard.writeText(meetingID); }}>
                    <ContentCopyIcon />
                </IconButton>
            </MeetingID>
            <IconButton color='inherit' onClick={OnToggleMute}>
                {muted ? <MutedIcon /> : <UnmutedIcon />}
            </IconButton>
            {dangerSignOn && <WatchingIcon />}
        </div>
    )
}

export default MeetingFooter