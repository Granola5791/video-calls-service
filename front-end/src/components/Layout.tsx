import HamburgerMenu from './HamburgerMenu'
import { MenuOptions } from '../constants/hebrew-constants'
import { LogOut } from '../utils/logout'
import { Outlet } from 'react-router-dom'
import { LocalStorage } from '../constants/general-contants'
import { useNavigation } from '../utils/navigation'
import { IsAdmin } from '../utils/roles'

const Layout = () => {
    const {
        goToTranscripts,
    } = useNavigation();
    const role = localStorage.getItem(LocalStorage.role);
    const username = localStorage.getItem(LocalStorage.username);
    const adminOptions = [{ text: MenuOptions.admin.transcripts, onClick: goToTranscripts }];
    const userOptions = [] as { text: string, onClick: () => void }[];
    const logoutOption = [{ text: MenuOptions.disconnect, onClick: LogOut }];
    const options = userOptions.concat(IsAdmin(role) ? adminOptions : []);
    return (
        <>
            <HamburgerMenu
                title={MenuOptions.title + username}
                topButtons={options}
                bottomButtons={logoutOption}
            />
            <main>
                <Outlet />
            </main>
        </>
    )
}

export default Layout