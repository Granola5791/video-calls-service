import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import LandingPage from './pages/LandingPage.tsx'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import LoginPage from './pages/LoginPage.tsx'
import SignupPage from './pages/SignupPage.tsx'
import HomePage from './pages/HomePage.tsx'
import { CheckAdminLoader, CheckLoginLoader } from './utils/check-login.ts'
import TestPage from './pages/TestPage.tsx'
import MeetingPage from './pages/MeetingPage.tsx'
import MeetingInfoPage from './pages/MeetingInfoPage.tsx'
import Layout from './components/Layout.tsx'
import { RouterPaths } from './constants/general-contants.ts'
import { CacheProvider } from '@emotion/react';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { cacheRtl } from './theme/cache';
import { theme } from './theme/theme';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import 'dayjs/locale/he';



const router = createBrowserRouter([
    {
        path: "/",
        element: <Layout />,
        children: [
            {
                index: true,
                element: <HomePage />,
                loader: CheckLoginLoader,
            },
            {
                path: RouterPaths.meetingInfo,
                element: <MeetingInfoPage />,
                loader: CheckAdminLoader,
            },
        ]
    },
    {
        path: RouterPaths.landing,
        element: <LandingPage />
    },
    {
        path: RouterPaths.login,
        element: <LoginPage />
    },
    {
        path: RouterPaths.signup,
        element: <SignupPage />
    },
    {
        path: RouterPaths.meeting,
        element: <MeetingPage />,
        loader: CheckLoginLoader,
    },
    {
        path: "/test",
        element: <TestPage />,
    },
])

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <CacheProvider value={cacheRtl}>
            <ThemeProvider theme={theme}>
                <LocalizationProvider dateAdapter={AdapterDayjs} adapterLocale='he'>
                    <CssBaseline />
                    <div dir="rtl" style={{ minHeight: '100vh' }}>
                        <RouterProvider router={router} />
                    </div>
                </LocalizationProvider>
            </ThemeProvider>
        </CacheProvider>
    </StrictMode>,
)
