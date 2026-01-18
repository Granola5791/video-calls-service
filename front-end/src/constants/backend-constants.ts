export const BackendAddressHttp = "http://localhost:8081";
export const BackendAddressWS = "ws://localhost:8081";
export const DasherServerAddressWS = "ws://localhost:8082";
export const DasherServerAddressHttp = "http://localhost:8082";

export const ApiEndpoints = {
    signUp: "/signup",
    logIn: "/login",
    checkLoginApi: "/check-login",
    checkAdminApi: "/check-admin",
    startStream: "/stream/{meeting_id}",
    getStream: "/get-stream/{meeting_id}/{user_id}/stream.mpd",
    createMeeting: "/create-meeting",
    getCallParticipants: "/get-call-participants",
    joinMeeting: "/join-meeting/{meeting_id}",
    getCallNotifications: "/get-call-notifications/{meeting_id}",
};

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