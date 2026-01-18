import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import LandingPage from './pages/LandingPage.tsx'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import LoginPage from './pages/LoginPage.tsx'
import SignupPage from './pages/SignupPage.tsx'
import HomePage from './pages/HomePage.tsx'
import { CheckLoginLoader } from './utils/check-login.ts'
import TestPage from './pages/TestPage.tsx'
import WebSocketWebCam from './components/WebSocketWebCam.tsx'
import { DasherServerAddressWS, ApiEndpoints } from './constants/backend-constants.ts'
import MeetingPage from './pages/MeetingPage.tsx'

const router = createBrowserRouter([
    {
        path: "/",
        element: <LandingPage />
    },
    {
        path: "/login",
        element: <LoginPage />
    },
    {
        path: "/signup",
        element: <SignupPage />
    },
    {
        path: "/home",
        element: <HomePage />,
        loader: CheckLoginLoader,
    },
    {
        path: "/meeting/:meetingID",
        element: <MeetingPage />,
        loader: CheckLoginLoader,
    },
    {
        path: "/test",
        element: <TestPage />
    },
    {
        path: "/test2",
        element: <WebSocketWebCam wsUrl={DasherServerAddressWS + ApiEndpoints.startStream} />
    },
])

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <RouterProvider router={router} />
    </StrictMode>,
)
