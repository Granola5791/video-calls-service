export const BackendAddress = "http://localhost:8081";
export const DasherServerAddressWS = "ws://localhost:8082";
export const DasherServerAddressHttp = "http://localhost:8082";

export const ApiEndpoints = {
    signUp: "/signup",
    logIn: "/login",
    checkLoginApi: "/check-login",
    checkAdminApi: "/check-admin",
    startStream: "/stream",
    getStream: "/get_stream/stream.mpd",
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