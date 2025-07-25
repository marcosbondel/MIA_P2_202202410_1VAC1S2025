import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { BrowserRouter } from 'react-router-dom'
import MIAApp from './MIAApp.jsx'
import { AppProvider } from './context/AppProvider.jsx'

createRoot(document.getElementById('root')).render(
    // <StrictMode>
        <BrowserRouter>
            <MIAApp />
        </BrowserRouter>
    // </StrictMode>
)
