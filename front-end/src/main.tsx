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
import MeetingTranscriptPage from './pages/MeetingTranscriptPage.tsx'
import Layout from './components/Layout.tsx'
import { RouterPaths } from './constants/general-contants.ts'
import SummaryPage from './pages/SummaryPage.tsx'

const router = createBrowserRouter([
    {
        path: "/",
        element: <Layout/>,
        children: [
            {
                index: true,
                element: <HomePage/>,
                loader: CheckLoginLoader,
            },
            {
                path: RouterPaths.meetingInfo,
                element: <MeetingInfoPage/>,
                loader: CheckAdminLoader,
            },
            {
                path: RouterPaths.meetingTranscript,
                element: <MeetingTranscriptPage/>,
                loader: CheckAdminLoader,
            },
            {
                path: RouterPaths.summary,
                element: <SummaryPage/>,
                loader: CheckAdminLoader,
            }
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
        <RouterProvider router={router} />
    </StrictMode>,
)
