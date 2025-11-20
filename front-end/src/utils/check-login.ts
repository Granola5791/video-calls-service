import { BackendAddress, ApiEndpoints, HttpStatuses, HttpStatusCodes } from "../constants/backend-constants";

async function fetchWithCredentials(apiPath: string) {
    return fetch(apiPath, {
        method: 'GET',
        credentials: 'include'
    });
}

export async function CheckLoginLoader() {
    const res = await fetchWithCredentials(BackendAddress + ApiEndpoints.checkAdminApi);
    if (!res.ok) {
        throw new Response(HttpStatuses.unauthorized, {
            status: HttpStatusCodes.found,
            headers: { Location: '/login' },
        });
    }
    return null;
}

export async function CheckAdminLoader() {
    const res = await fetchWithCredentials(BackendAddress + ApiEndpoints.checkAdminApi);
    if (!res.ok) {
        throw new Response(HttpStatuses.unauthorized, {
            status: HttpStatusCodes.found,
            headers: { Location: '/login' },
        });
    }
    return null;
}