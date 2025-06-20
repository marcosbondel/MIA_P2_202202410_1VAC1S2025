import React from 'react'
import { Route, Routes } from 'react-router-dom'
import { LoginPage } from '../auth/pages'
import { MIAPage } from '../mia/pages'

export const AppRouter = () => {
    return (
        <Routes>
            <Route path="/" element={<LoginPage/>}/>
            <Route path="/mia" element={<MIAPage/>}/>
            <Route path="/*" element={<LoginPage/>}/>
        </Routes>
    )
}
