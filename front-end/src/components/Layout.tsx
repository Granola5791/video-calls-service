import HamburgerMenu from './HamburgerMenu'
import { MenuOptions } from '../constants/hebrew-constants'
import { LogOut } from '../utils/logout'
import { Outlet } from 'react-router-dom'
import { LocalStorage } from '../constants/general-contants'
import { useNavigation } from '../utils/navigation'
import { IsAdmin } from '../utils/roles'
import { MediumScreen } from '../styled-components/StyledBoxes'
import type { MenuOption } from '../types/menuOptions'

const Layout = () => {
    const {
        goToMeetingInfo: goToTranscripts,
        goToHome,
    } = useNavigation();
    const role = localStorage.getItem(LocalStorage.role);
    const username = localStorage.getItem(LocalStorage.username);
    const adminOptions = [{ text: MenuOptions.admin.meetingInfos, onClick: goToTranscripts }];
    const userOptions = [] as MenuOption[];
    const logoutOption = [{ text: MenuOptions.disconnect, onClick: LogOut }] as MenuOption[];
    const options = userOptions.concat(IsAdmin(role) ? adminOptions : []);
    return (
        <>
            <HamburgerMenu
                title={MenuOptions.title + username}
                topButtons={options}
                bottomButtons={logoutOption}
            />

            <main>
                <MediumScreen>
                    <Outlet />
                </MediumScreen>
            </main>
        </>
    )
}

export default Layout