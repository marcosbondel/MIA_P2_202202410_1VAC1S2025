import React from 'react'
import { Route, Routes } from 'react-router-dom'
import { LoginPage } from '../auth/pages'

export const AppRouter = () => {
    return (
        <Routes>
            <Route path="/" element={<LoginPage/>}/>
        </Routes>
    )
}
