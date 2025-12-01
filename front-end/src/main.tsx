import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import LandingPage from './pages/LandingPage.tsx'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import LoginPage from './pages/LoginPage.tsx'
import SignupPage from './pages/SignupPage.tsx'
import HomePage from './pages/HomePage.tsx'
import { CheckLoginLoader } from './utils/check-login.ts'
import WebCamChunks from './components/WebCamChunks.tsx'


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
        path: "/test",
        element: <WebCamChunks />
    }
])

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <RouterProvider router={router} />
    </StrictMode>,
)
