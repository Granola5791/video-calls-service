import { BackendAddressHttp, ApiEndpoints, HttpStatuses, HttpStatusCodes } from "../constants/backend-constants";

async function fetchWithCredentials(apiPath: string) {
    const res = await fetch(apiPath, {
        method: 'GET',
        credentials: 'include',
    });
    return res;
}

export async function CheckLoginLoader() {
    const res = await fetchWithCredentials(BackendAddressHttp + ApiEndpoints.checkLoginApi);
    if (!res.ok) {
        throw new Response(HttpStatuses.unauthorized, {
            status: HttpStatusCodes.found,
            headers: { Location: '/login' },
        });
    }
    return null;
}

export async function CheckAdminLoader() {
    const res = await fetchWithCredentials(BackendAddressHttp + ApiEndpoints.checkAdminApi);
    if (!res.ok) {
        throw new Response(HttpStatuses.unauthorized, {
            status: HttpStatusCodes.found,
            headers: { Location: '/login' },
        });
    }
    return null;
}