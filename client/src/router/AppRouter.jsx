import React from 'react'
import { Route, Routes } from 'react-router-dom'
import { LoginPage } from '../auth/pages'
import { MIAPage } from '../mia/pages'
import { PublicRoute } from './PublicRoute'
import { PrivateRoute } from './PrivateRoute'

export const AppRouter = () => {
    return (
        <Routes>
            <Route path="/login" element={
                <PublicRoute>
                    <LoginPage/>
                </PublicRoute>
            }/>
            <Route path="/*" element={
                <PrivateRoute>
                    <MIAPage/>
                </PrivateRoute>
            }/>

        </Routes>
    )
}
