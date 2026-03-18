import HamburgerMenu from './HamburgerMenu'
import { MenuOptions } from '../constants/hebrew-constants'
import { LogOut } from '../utils/logout'
import { Outlet } from 'react-router-dom'

const Layout = () => {
    return (
        <>
            <HamburgerMenu
                topButtons={[{ text: MenuOptions.disconnect, onClick: LogOut }]}
            />
            <main>
                <Outlet />
            </main>
        </>
    )
}

export default Layout