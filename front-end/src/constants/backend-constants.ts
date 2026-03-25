import type { QueryParam } from "../types/queryParam";

export const DasherServer = {
    httpAddress: "https://dasherserver.local.my:8082",
    wsAddress: "wss://dasherserver.local.my:8082",
    api: {
        startStream: "/stream/{meeting_id}",
        getStream: "/get-stream/{meeting_id}/{user_id}/stream.mpd",
        createMeeting: "/create-meeting",
        joinMeeting: "/join-meeting/{meeting_id}",
    },
}

export const UsersServer = {
    httpAddress: "https://usersserver.local.my:8081",
    wsAddress: "wss://usersserver.local.my:8081",
    api: {
        queryParams: {
            from: "from",
            to: "to",
            host_username: "host_username",
        },
        signUp: "/signup",
        logIn: "/login",
        logOut: "/logout",
        checkLoginApi: "/check-login",
        checkAdminApi: "/check-admin",
        createMeeting: "/create-meeting/{require_face}",
        getCallParticipants: "/get-call-participants",
        joinMeeting: "/join-meeting/{meeting_id}",
        getCallNotifications: "/get-call-notifications/{meeting_id}",
        leaveMeeting: "/leave-meeting/{meeting_id}",
        keepAlive: "/keep-alive/{meeting_id}",
        kickParticipant: "/kick-participant/{meeting_id}/{participant_id}",
        isAbleToJoinMeeting: "/is-able-to-join-meeting/{meeting_id}",
        getMeetingInfos: "/get-meeting-infos",
        getTranscript: "/get-transcript/{meeting_id}/{participant_id}",
        getAllMeetingParticipants: "/get-all-meeting-participants/{meeting_id}",
        getSummary: "/get-summary/{meeting_id}",
    },
}

export const HttpStatusCodes = {
    OK: 200,
    Created: 201,
    found: 302,
    BadRequest: 400,
    Unauthorized: 401,
    Forbidden: 403,
    NotFound: 404,
    Conflict: 409,
    InternalServerError: 500,
};

export const HttpStatuses = {
    ok: "OK",
    created: "Created",
    badRequest: "Bad Request",
    unauthorized: "Unauthorized",
    forbidden: "Forbidden",
    notFound: "Not Found",
    conflict: "Conflict",
    internalServerError: "Internal Server Error",
};

export const CallEventTypes = {
    participantJoined: 0,
    participantLeft: 1,
    participantKickedByHost: 2,
    meetingEnded: 3,
}

export const SetUrlParams = (url: string, ...params: any[]): string => {
    let result = url;
    let paramIndex = 0;

    result = result.replace(/\{[^}]+\}/g, () => {
        if (paramIndex < params.length) {
            return String(params[paramIndex++]);
        }
        return '';
    });

    return result;
}

export const AddQueryParams = (url: string, params: QueryParam[]): string => {
    let result = url + '?';
    for (let i = 0; i < params.length; i++) {
        result += `${params[i].key}=${params[i].value}&`;
    }
    return result;
}