export type TranscriptSegment = {
    text: string,
    start: number,
    end: number
}

export type UserTranscript = {
    username: string,
    segments: TranscriptSegment[]
}

export type Index2D = {
    x: number,
    y: number
}