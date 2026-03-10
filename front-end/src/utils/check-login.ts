import { UsersServer, HttpStatuses, HttpStatusCodes } from "../constants/backend-constants";

async function fetchWithCredentials(apiPath: string) {
    const res = await fetch(apiPath, {
        method: 'GET',
        credentials: 'include',
    });
    return res;
}

export async function CheckLoginLoader() {
    const res = await fetchWithCredentials(UsersServer.httpAddress + UsersServer.api.checkLoginApi);
    if (!res.ok) {
        throw new Response(HttpStatuses.unauthorized, {
            status: HttpStatusCodes.found,
            headers: { Location: '/' },
        });
    }
    return null;
}

export async function CheckAdminLoader() {
    const res = await fetchWithCredentials(UsersServer.httpAddress + UsersServer.api.checkAdminApi);
    if (!res.ok) {
        throw new Response(HttpStatuses.unauthorized, {
            status: HttpStatusCodes.found,
            headers: { Location: '/' },
        });
    }
    return null;
}