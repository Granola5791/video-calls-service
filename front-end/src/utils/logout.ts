import { UsersServer } from "../constants/backend-constants";
import { RouterPaths } from "../constants/general-contants";

export const LogOut = () => {
    const LogoutFrontend = () => {
        localStorage.clear();
        window.location.href = RouterPaths.landing;
    }

    const LogoutBackend = async () => {
        await fetch(UsersServer.httpAddress + UsersServer.api.logOut, {
            method: 'POST',
            credentials: 'include',
        });
    }

    LogoutBackend().then(LogoutFrontend);
}