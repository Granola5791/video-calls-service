export const Colors = {
    primary: '#ff1a8dff',
    primaryHover: '#ff52b4ff',
    meetingBackground: '#111',
    danger: 'rgb(255, 0, 0)',
    dangerHover: 'rgb(255, 60, 60)',
};

export const AuthRequirements = {
    passwordMinLength: 10,
    passwordMaxLength: 64,
    usernameMinLength: 1,
    usernameMaxLength: 32,
}

export const WebCamConfig = {
    defaultWidth: '640px',
    defaultHeight: '480px',
}

export const StreamConfig = {
    serverReadyMsg: 'READY',
    chunkIntervalMs: 2000, // 2 seconds
    sentChunkMsg: 'Sent chunk at',
    checkStreamAvailabilityIntervalMs: 2000, // 2 seconds
}

export const MeetingConfig = {
    exitWaitTimeMs: 5000,
    meetingState: {
        none: 0,
        wrongID: 1,
        active: 2,
        ended: 3,
        left: 4,
        kicked: 5,
    }
}