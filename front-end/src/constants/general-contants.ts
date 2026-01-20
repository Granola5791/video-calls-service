export const Colors = {
    primary: '#ff1a8dff',
    primaryHover: '#ff52b4ff',
    meetingBackground: '#111',
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