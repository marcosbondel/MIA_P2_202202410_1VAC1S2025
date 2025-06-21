import React from 'react'
import { Route, Routes } from 'react-router-dom'
import { LoginPage } from '../auth/pages'
import { MIAPage, TerminalPage } from '../mia/pages'
import { PublicRoute } from './PublicRoute'
import { PrivateRoute } from './PrivateRoute'

export const AppRouter = () => {
    return (
        <Routes>
            {/* Public Route */}
            <Route path="/login" element={
                <PublicRoute>
                    <LoginPage />
                </PublicRoute>
            } />

            {/* Private Routes */}
            <Route path="/mia" element={
                <PrivateRoute>
                    <MIAPage />
                </PrivateRoute>
            } />

            <Route path="/mia/terminal" element={
                <PrivateRoute>
                    <TerminalPage />
                </PrivateRoute>
            } />
            
            <Route path="/*" element={
                <PrivateRoute>
                    <MIAPage/>
                </PrivateRoute>
            }/>
        </Routes>
    )
}

